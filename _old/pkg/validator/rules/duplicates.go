package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Ensure unique gene identifiers within the composed genome.
func init() {
	core.Register(checkDuplicateGenes)
}

func checkDuplicateGenes(genome *loader.ComposedGenome, res *core.Result) {
	index := map[string]int{}
	for i, g := range genome.Genes {
		key := g.Chromosome + "." + g.Name
		if prev, exists := index[key]; exists {
			res.AddWithSeverity(severityFor(core.CategoryUniqueness, core.Error), fmt.Sprintf("genes[%d]", i), fmt.Sprintf("duplicate gene identifier %q (also at genes[%d])", key, prev))
		} else {
			index[key] = i
		}
	}
}
