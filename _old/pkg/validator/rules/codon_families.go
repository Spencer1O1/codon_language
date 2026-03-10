package rules

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkCodonFamilies)
}

var coreBuckets = map[string]struct{}{
	"entities":     {},
	"capabilities": {},
	"relations":    {},
	"references":   {},
	"traits":       {},
}

// checkCodonFamilies ensures non-core codon buckets are declared in the codon family registry when present.
func checkCodonFamilies(genome *loader.ComposedGenome, res *core.Result) {
	if len(genome.CodonFamilies) == 0 {
		return // no registry provided; skip
	}
	for gi, g := range genome.Genes {
		for bucket := range g.RawCodons {
			if _, ok := coreBuckets[bucket]; ok {
				continue
			}
			if _, ok := genome.CodonFamilies[bucket]; !ok {
				res.AddWithSeverity(severityFor(core.CategoryStructural, core.Error), genePath(gi), "codon bucket \""+bucket+"\" has no declared codon family")
			}
		}
	}
}
