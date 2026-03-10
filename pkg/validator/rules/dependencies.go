package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Ensure declared dependencies exist as genes in the composed genome.
func init() {
	core.Register(checkDependenciesExist)
}

func checkDependenciesExist(genome *loader.ComposedGenome, res *core.Result) {
	geneIndex := map[string]struct{}{}
	for _, g := range genome.Genes {
		geneIndex[g.Chromosome+"."+g.Name] = struct{}{}
	}
	for _, g := range genome.Genes {
		for _, dep := range g.Dependencies {
			if _, ok := geneIndex[dep]; !ok {
				res.Add(geneKey(g), fmt.Sprintf("dependency %q not found in genome", dep))
			}
		}
	}
}

func geneKey(g loader.ComposedGene) string {
	return fmt.Sprintf("gene[%s.%s]", g.Chromosome, g.Name)
}
