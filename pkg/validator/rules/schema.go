package rules

import (
	"fmt"
	"regexp"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(schemaRules)
}

// schemaRules validates codon instances against their codon schemas using TypeExpr.
// This is a shallow validator (object/map/list/scalar, Optional, Union, Map<K,V>, Object literals, Regex).
func schemaRules(g *loader.Genome, env map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for codonName, val := range gene.Codons {
			schema, ok := g.Schemas[codonName]
			if !ok {
				// basicShape already reports missing schemas
				continue
			}
			if err := validateValueAgainstType(val, schema.TypeAST, env); err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_mismatch", Message: fmt.Sprintf("%v", err), Gene: gene.Name, Codon: codonName})
			}
		}
	}
}

func validateValueAgainstType(v any, t nt.TypeNode, env map[string]nt.TypeNode) error {
	switch node := t.(type) {
	case nt.OptionalType:
		if v == nil {
			return nil
		}
		return validateValueAgainstType(v, node.Base, env)
	case nt.ListType:
		arr, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected list")
		}
		for _, elem := range arr {
			if err := validateValueAgainstType(elem, node.Base, env); err != nil {
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
			// check key and value types when provided
			if len(node.Args) == 2 {
				keyType := node.Args[0]
				valType := node.Args[1]
				for k, vv := range m {
					if err := validateMapKey(k, keyType, env); err != nil {
						return fmt.Errorf("key %q: %v", k, err)
					}
					if err := validateValueAgainstType(vv, valType, env); err != nil {
						return err
					}
				}
			}
			return nil
		}
		if node.Name == "Regex" {
			pat := extractRegexPattern(node.Args)
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected string matching regex %q", pat)
			}
			re, err := regexp.Compile(pat)
			if err != nil {
				return fmt.Errorf("invalid regex pattern %q", pat)
			}
			if !re.MatchString(s) {
				return fmt.Errorf("value %q does not match regex %q", s, pat)
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
				if err := validateValueAgainstType(vv, f.Type, env); err != nil {
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
			if err := validateValueAgainstType(v, opt, env); err == nil {
				return nil
			}
		}
		return fmt.Errorf("value did not match any union option")
	case nt.NameType:
		switch node.Name {
		case "string", "number", "boolean", "uuid", "datetime", "json", "yaml":
			if isScalar(v) {
				return nil
			}
			return fmt.Errorf("expected scalar for %s", node.Name)
		case "ref":
			if s, ok := v.(string); ok && refPattern.MatchString(s) {
				return nil
			}
			return fmt.Errorf("expected ref path string")
		case "TypeExpr":
			if isScalar(v) {
				return nil
			}
			return fmt.Errorf("expected scalar for %s", node.Name)
		default:
			if resolved, ok := env[node.Name]; ok {
				return validateValueAgainstType(v, resolved, env)
			}
			return fmt.Errorf("unknown type %s", node.Name)
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

func validateMapKey(key string, t nt.TypeNode, env map[string]nt.TypeNode) error {
	switch kt := t.(type) {
	case nt.NameType:
		// primitives as map keys: accept string form
		switch kt.Name {
		case "string", "number", "boolean", "uuid", "datetime", "json", "yaml", "ref", "TypeExpr", "primitive", "any":
			return nil
		}
		if resolved, ok := env[kt.Name]; ok {
			return validateMapKey(key, resolved, env)
		}
		// default to string match
		return nil
	case nt.GenericType:
		if kt.Name == "Regex" {
			pat := extractRegexPattern(kt.Args)
			re, err := regexp.Compile(pat)
			if err != nil {
				return fmt.Errorf("invalid regex pattern %q", pat)
			}
			if !re.MatchString(key) {
				return fmt.Errorf("key does not match regex %q", pat)
			}
			return nil
		}
	}
	return nil
}

func extractRegexPattern(args []nt.TypeNode) string {
	if len(args) == 0 {
		return ""
	}
	if lit, ok := args[0].(nt.LiteralType); ok {
		return lit.Value
	}
	return ""
}

var refPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*(\.[A-Za-z][A-Za-z0-9_]*){0,3}$`)
