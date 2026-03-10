package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkRelationsAndReferences)
}

func checkRelationsAndReferences(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		entities := map[string]struct{}{}
		for _, e := range g.Entities {
			entities[e.Name] = struct{}{}
		}
		for ri, r := range g.Relations {
			if _, ok := entities[r.From]; !ok {
				res.Add(fmt.Sprintf("genes[%d].relations[%d]", gi, ri), fmt.Sprintf("from entity %q not defined in gene", r.From))
			}
			if _, ok := entities[r.To]; !ok {
				res.Add(fmt.Sprintf("genes[%d].relations[%d]", gi, ri), fmt.Sprintf("to entity %q not defined in gene", r.To))
			}
		}
		for ri, r := range g.References {
			if _, ok := entities[r.From]; !ok {
				res.Add(fmt.Sprintf("genes[%d].references[%d]", gi, ri), fmt.Sprintf("from entity %q not defined in gene", r.From))
			}
		}
	}
}
