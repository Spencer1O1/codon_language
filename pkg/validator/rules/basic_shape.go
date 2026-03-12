package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(basicShape)
}

// basicShape checks that each codon has a codon schema and matches top-level shape.
func basicShape(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		// filename stem must match gene name
		stem := strings.TrimSuffix(filepath.Base(gene.Path), ".yaml")
		if stem != "" && stem != gene.Name {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "gene_filename_matches_name", Message: fmt.Sprintf("gene filename %s must match gene name %s", stem, gene.Name), Gene: gene.Name})
		}
		for bucket, val := range gene.Codons {
			schema, ok := g.Schemas[bucket]
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_missing", Message: fmt.Sprintf("codon %s has no codon schema definition", bucket), Gene: gene.Name, Codon: bucket})
				continue
			}
			switch ft := schema.TypeAST.(type) {
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
