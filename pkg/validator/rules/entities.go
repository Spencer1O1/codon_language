package rules

import (
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.Register(entityRules) }

func entityRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		codon, ok := gene.Codons["entities"].(map[string]any)
		if !ok {
			continue
		}
		for name, raw := range codon {
			if strings.TrimSpace(name) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "entity_key_required", Message: "entity names must be non-empty", Gene: gene.Name, Codon: "entities"})
				continue
			}
			fields, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			for fname, fval := range fields {
				if strings.TrimSpace(fname) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "field_key_format", Message: "field names must be lower_snake_case", Gene: gene.Name, Codon: "entities"})
				}
				// field type required: value is scalar (interpreted as type_expr) OK; if map, must include type/type_expr/ref
				if m, ok := fval.(map[string]any); ok {
					_, hasType := m["type"]
					_, hasTypeExpr := m["type_expr"]
					_, hasRef := m["ref"]
					if !hasType && !hasTypeExpr && !hasRef {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "field_type_required", Message: "field must declare type/type_expr/ref", Gene: gene.Name, Codon: "entities"})
					}
				} else {
					// scalar: treated as type_expr, okay
				}
			}
		}
	}
}
