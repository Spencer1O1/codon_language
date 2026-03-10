package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Enforce unique relation/reference names within a gene.
func init() {
	core.Register(checkRelationReferenceNames)
	core.Register(checkEntityCapabilityNames)
}

func checkRelationReferenceNames(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		relNames := map[string]int{}
		for i, r := range g.Relations {
			if prev, ok := relNames[r.Name]; ok {
				res.AddWithSeverity(core.Error, fmt.Sprintf("%s.relations[%d]", genePath(gi), i), fmt.Sprintf("duplicate relation name %q (also at relations[%d])", r.Name, prev))
			} else {
				relNames[r.Name] = i
			}
		}
		refNames := map[string]int{}
		for i, r := range g.References {
			if prev, ok := refNames[r.Name]; ok {
				res.AddWithSeverity(core.Error, fmt.Sprintf("%s.references[%d]", genePath(gi), i), fmt.Sprintf("duplicate reference name %q (also at references[%d])", r.Name, prev))
			} else {
				refNames[r.Name] = i
			}
		}
	}
}

// Enforce unique entity and capability names within a gene.
func checkEntityCapabilityNames(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		entityNames := map[string]int{}
		for i, e := range g.Entities {
			if prev, ok := entityNames[e.Name]; ok {
				res.AddWithSeverity(severityFor(core.CategoryUniqueness, core.Error),
					genePath(gi), fmt.Sprintf("duplicate entity name %q (also at entities[%d])", e.Name, prev))
			} else {
				entityNames[e.Name] = i
			}
		}
		capNames := map[string]int{}
		for i, c := range g.Capabilities {
			if prev, ok := capNames[c.Name]; ok {
				res.AddWithSeverity(severityFor(core.CategoryUniqueness, core.Error),
					genePath(gi), fmt.Sprintf("duplicate capability name %q (also at capabilities[%d])", c.Name, prev))
			} else {
				capNames[c.Name] = i
			}
		}
	}
}
