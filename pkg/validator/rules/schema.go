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
			// Skip schema validation for relations to avoid recursive object types; other rules cover them.
			if codonName == "relations" {
				continue
			}
			schema, ok := g.Schemas[codonName]
			if !ok {
				// basicShape already reports missing schemas
				continue
			}
			if issue := validateValueAgainstType(val, schema.TypeAST, env, g, gene, codonName, res, 0); issue != nil {
				res.Add(*issue)
			}
		}
	}
}

// validateValueAgainstType returns a coded issue on the first failure, or nil on success.
func validateValueAgainstType(v any, t nt.TypeNode, env map[string]nt.TypeNode, genome *loader.Genome, gene loader.Gene, codon string, res *core.Result, depth int) *core.Issue {
	if depth > 256 {
		return issue("schema_mismatch", "schema recursion too deep", gene.Name, codon)
	}
	switch node := t.(type) {
	case nt.OptionalType:
		if v == nil {
			return nil
		}
		return validateValueAgainstType(v, node.Base, env, genome, gene, codon, res, depth+1)
	case nt.ListType:
		arr, ok := v.([]any)
		if !ok {
			return issue("schema_mismatch", "expected list", gene.Name, codon)
		}
		for _, elem := range arr {
			if err := validateValueAgainstType(elem, node.Base, env, genome, gene, codon, res, depth+1); err != nil {
				return err
			}
		}
		return nil
	case nt.GenericType:
		if node.Name == "Map" {
			m, ok := v.(map[string]any)
			if !ok {
				return issue("schema_mismatch", "expected map", gene.Name, codon)
			}
			// check key and value types when provided
			if len(node.Args) == 2 {
				keyType := node.Args[0]
				valType := node.Args[1]
				for k, vv := range m {
					if err := validateMapKey(k, keyType, env); err != nil {
						return issue("map_key_constraint", fmt.Sprintf("key %q: %v", k, err), gene.Name, codon)
					}
					if err := validateValueAgainstType(vv, valType, env, genome, gene, codon, res, depth+1); err != nil {
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
				return issue("regex_constraint_violation", fmt.Sprintf("expected string matching regex %q", pat), gene.Name, codon)
			}
			re, err := regexp.Compile(pat)
			if err != nil {
				return issue("regex_constraint_violation", fmt.Sprintf("invalid regex pattern %q", pat), gene.Name, codon)
			}
			if !re.MatchString(s) {
				return issue("regex_constraint_violation", fmt.Sprintf("value %q does not match regex %q", s, pat), gene.Name, codon)
			}
			return nil
		}
		// other generics treated as scalar here
		return nil
	case nt.ObjectType:
		m, ok := v.(map[string]any)
		if !ok {
			return issue("schema_mismatch", "expected object", gene.Name, codon)
		}
		// required/optional via presence; no required flag in AST, but we can ensure defined fields exist
		for _, f := range node.Fields {
			if vv, ok := m[f.Name]; ok {
				if err := validateValueAgainstType(vv, f.Type, env, genome, gene, codon, res, depth+1); err != nil {
					return err
				}
			} else {
				// treat missing as error unless OptionalType
				if _, isOpt := f.Type.(nt.OptionalType); !isOpt {
					return issue("schema_mismatch", fmt.Sprintf("field %s required", f.Name), gene.Name, codon)
				}
			}
		}
		return nil
	case nt.UnionType:
		// accept if any branch matches
		for _, opt := range node.Options {
			if err := validateValueAgainstType(v, opt, env, genome, gene, codon, res, depth+1); err == nil {
				return nil
			}
		}
		return issue("schema_mismatch", "value did not match any union option", gene.Name, codon)
	case nt.NameType:
		switch node.Name {
		case "string", "number", "boolean", "uuid", "datetime", "json", "yaml":
			if isScalar(v) {
				return nil
			}
			return issue("schema_mismatch", fmt.Sprintf("expected scalar for %s", node.Name), gene.Name, codon)
		case "object":
			// Treat generic object as an escape hatch; accept maps without deep recursion to avoid self-referential loops.
			if _, ok := v.(map[string]any); ok {
				return nil
			}
			return issue("schema_mismatch", "expected object", gene.Name, codon)
		case "ref":
			if s, ok := v.(string); ok {
				// Reuse reference resolution so ref TypeExpr behaves like {ref: ...}
				checkRef(s, genome, gene, codon, res)
				return nil
			}
			return issue("schema_mismatch", "expected ref path string", gene.Name, codon)
		case "TypeExpr":
			if isScalar(v) {
				return nil
			}
			return issue("schema_mismatch", fmt.Sprintf("expected scalar for %s", node.Name), gene.Name, codon)
		default:
			if resolved, ok := env[node.Name]; ok {
				return validateValueAgainstType(v, resolved, env, genome, gene, codon, res, depth+1)
			}
			return issue("type_name_unknown", fmt.Sprintf("unknown type %s", node.Name), gene.Name, codon)
		}
	case nt.LiteralType:
		if s, ok := v.(string); ok && s == node.Value {
			return nil
		}
		return issue("schema_mismatch", fmt.Sprintf("expected literal %s", node.Value), gene.Name, codon)
	default:
		// primitives/scalars
		if isScalar(v) {
			return nil
		}
		return issue("schema_mismatch", "expected scalar", gene.Name, codon)
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
		return fmt.Errorf("unknown type %s", kt.Name)
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

func issue(code, msg, gene, codon string) *core.Issue {
	return &core.Issue{Severity: core.SeverityError, Code: code, Message: msg, Gene: gene, Codon: codon}
}
