package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkTraitDuplicates)
}

func checkTraitDuplicates(genome *loader.ComposedGenome, res *core.Result) {
	seen := map[string]int{}
	for i, t := range genome.Traits {
		if prev, ok := seen[t]; ok {
			res.AddWithSeverity(core.Warning, "genome.traits", fmt.Sprintf("duplicate trait %q (indexes %d and %d)", t, prev, i))
		} else {
			seen[t] = i
		}
	}
	for gi, g := range genome.Genes {
		seenGene := map[string]int{}
		for i, t := range g.Traits {
			if prev, ok := seenGene[t]; ok {
				res.AddWithSeverity(core.Warning, genePath(gi)+".traits", fmt.Sprintf("duplicate trait %q (indexes %d and %d)", t, prev, i))
			} else {
				seenGene[t] = i
			}
		}
	}
}
