package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Enforce unique relation/reference names within a gene.
func init() {
	core.Register(checkRelationReferenceNames)
}

func checkRelationReferenceNames(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		relNames := map[string]int{}
		for i, r := range g.Relations {
			if prev, ok := relNames[r.Name]; ok {
				res.Add(fmt.Sprintf("%s.relations[%d]", genePath(gi), i), fmt.Sprintf("duplicate relation name %q (also at relations[%d])", r.Name, prev))
			} else {
				relNames[r.Name] = i
			}
		}
		refNames := map[string]int{}
		for i, r := range g.References {
			if prev, ok := refNames[r.Name]; ok {
				res.Add(fmt.Sprintf("%s.references[%d]", genePath(gi), i), fmt.Sprintf("duplicate reference name %q (also at references[%d])", r.Name, prev))
			} else {
				refNames[r.Name] = i
			}
		}
	}
}
