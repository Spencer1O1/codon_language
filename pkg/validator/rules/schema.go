package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(schemaRules)
}

// schemaRules validates codon instances against their codon schemas using TypeExpr.
// This is a shallow validator (object/map/list/scalar, Optional, Union, Map<K,V>, Object literals).
func schemaRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for codonName, val := range gene.Codons {
			fam, ok := g.Families[codonName]
			if !ok {
				// basicShape already reports missing families
				continue
			}
			if err := validateValueAgainstType(val, fam.TypeAST); err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_mismatch", Message: fmt.Sprintf("%v", err), Gene: gene.Name, Codon: codonName})
			}
		}
	}
}

func validateValueAgainstType(v any, t nt.TypeNode) error {
	switch node := t.(type) {
	case nt.OptionalType:
		if v == nil {
			return nil
		}
		return validateValueAgainstType(v, node.Base)
	case nt.ListType:
		arr, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected list")
		}
		for _, elem := range arr {
			if err := validateValueAgainstType(elem, node.Base); err != nil {
				return err
			}
		}
		return nil
	case nt.GenericType:
		if node.Name == "Map" {
			m, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("expected map")
			}
			// only check value type (keys are strings in YAML)
			if len(node.Args) == 2 {
				valType := node.Args[1]
				for _, vv := range m {
					if err := validateValueAgainstType(vv, valType); err != nil {
						return err
					}
				}
			}
			return nil
		}
		// other generics treated as scalar here
		return nil
	case nt.ObjectType:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		// required/optional via presence; no required flag in AST, but we can ensure defined fields exist
		for _, f := range node.Fields {
			if vv, ok := m[f.Name]; ok {
				if err := validateValueAgainstType(vv, f.Type); err != nil {
					return fmt.Errorf("field %s: %v", f.Name, err)
				}
			} else {
				// treat missing as error unless OptionalType
				if _, isOpt := f.Type.(nt.OptionalType); !isOpt {
					return fmt.Errorf("field %s required", f.Name)
				}
			}
		}
		return nil
	case nt.UnionType:
		// accept if any branch matches
		for _, opt := range node.Options {
			if err := validateValueAgainstType(v, opt); err == nil {
				return nil
			}
		}
		return fmt.Errorf("value did not match any union option")
	case nt.NameType:
		switch node.Name {
		case "string", "number", "boolean", "uuid", "datetime", "json", "yaml", "ref", "TypeExpr":
			if isScalar(v) {
				return nil
			}
			return fmt.Errorf("expected scalar for %s", node.Name)
		default:
			// unknown/other named types accepted as scalar here; could resolve against env
			return nil
		}
	case nt.LiteralType:
		if s, ok := v.(string); ok && s == node.Value {
			return nil
		}
		return fmt.Errorf("expected literal %s", node.Value)
	default:
		// primitives/scalars
		if isScalar(v) {
			return nil
		}
		return fmt.Errorf("expected scalar")
	}
}

func isScalar(v any) bool {
	switch v.(type) {
	case string, bool, float64, int, int64, nil:
		return true
	default:
		return false
	}
}
