package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(basicShape)
}

// basicShape checks that each codon has a family and matches top-level shape.
func basicShape(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for bucket, val := range gene.Codons {
			fam, ok := g.Families[bucket]
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "family_missing", Message: fmt.Sprintf("codon %s has no family definition", bucket), Gene: gene.Name, Codon: bucket})
				continue
			}
			switch ft := fam.TypeAST.(type) {
			case nt.ObjectType:
				if _, ok := val.(map[string]any); !ok {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "shape_mismatch", Message: "expected object", Gene: gene.Name, Codon: bucket})
				}
			case nt.GenericType:
				if ft.Name == "Map" {
					if _, ok := val.(map[string]any); !ok {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "shape_mismatch", Message: "expected map", Gene: gene.Name, Codon: bucket})
					}
				}
				if ft.Name == "List" {
					if _, ok := val.([]any); !ok {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "shape_mismatch", Message: "expected list", Gene: gene.Name, Codon: bucket})
					}
				}
			default:
				_ = val // accept scalar/union etc.; deeper checks belong elsewhere
			}
		}
	}
}
