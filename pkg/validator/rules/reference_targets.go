package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkReferenceTargets)
}

// checkReferenceTargets ensures all cross-gene references (reference codons,
// reference-type fields, capability I/O references) point to existing entities.
func checkReferenceTargets(genome *loader.ComposedGenome, res *core.Result) {
	entityIndex := make(map[string]struct{})
	for _, g := range genome.Genes {
		for _, e := range g.Entities {
			key := fmt.Sprintf("%s.%s.%s", g.Chromosome, g.Name, e.Name)
			entityIndex[key] = struct{}{}
		}
	}

	// reference codons
	for gi, g := range genome.Genes {
		for ri, r := range g.References {
			target := r.To
			if !exists(entityIndex, target) {
				res.AddWithSeverity(severityFor(core.CategoryReferences, core.Error), fmt.Sprintf("%s.references[%d]", genePath(gi), ri), fmt.Sprintf("reference target %q not found", target))
			}
		}

		// entity fields of type reference
		for _, e := range g.Entities {
			for fname, f := range e.Fields {
				if f.Type == "reference" && f.Reference != "" && !exists(entityIndex, f.Reference) {
					res.AddWithSeverity(severityFor(core.CategoryReferences, core.Error), fmt.Sprintf("%s.entities[%s].fields[%s]", genePath(gi), e.Name, fname), fmt.Sprintf("reference target %q not found", f.Reference))
				}
			}
		}

		// capability inputs/outputs of type reference
		for _, c := range g.Capabilities {
			for fname, f := range c.Inputs {
				if f.Type == "reference" && f.Reference != "" && !exists(entityIndex, f.Reference) {
					res.AddWithSeverity(severityFor(core.CategoryReferences, core.Error), fmt.Sprintf("%s.capabilities[%s].inputs[%s]", genePath(gi), c.Name, fname), fmt.Sprintf("reference target %q not found", f.Reference))
				}
			}
			for fname, f := range c.Outputs {
				if f.Type == "reference" && f.Reference != "" && !exists(entityIndex, f.Reference) {
					res.AddWithSeverity(severityFor(core.CategoryReferences, core.Error), fmt.Sprintf("%s.capabilities[%s].outputs[%s]", genePath(gi), c.Name, fname), fmt.Sprintf("reference target %q not found", f.Reference))
				}
			}
		}
	}
}

func exists(index map[string]struct{}, key string) bool {
	_, ok := index[key]
	return ok
}
