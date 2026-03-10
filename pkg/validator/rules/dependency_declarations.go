package rules

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkDependencyDeclarations)
}

// Ensure dependency lists are unique and cross-gene references declare deps.
func checkDependencyDeclarations(genome *loader.ComposedGenome, res *core.Result) {
	geneIndex := map[string]struct{}{}
	for _, g := range genome.Genes {
		geneIndex[g.Chromosome+"."+g.Name] = struct{}{}
	}

	for gi, g := range genome.Genes {
		seenDeps := map[string]int{}
		for di, dep := range g.Dependencies {
			if prev, exists := seenDeps[dep]; exists {
				res.Add(genePath(gi), fmt.Sprintf("duplicate dependency %q (also at index %d)", dep, prev))
			} else {
				seenDeps[dep] = di
			}
		}

		// references
		for ri, r := range g.References {
			targetGene, err := loader.GenePartFromEntityRef(r.To)
			if err != nil {
				res.Add(fmt.Sprintf("%s.references[%d]", genePath(gi), ri), err.Error())
				continue
			}
			if !declaresDependency(g.Dependencies, targetGene) && targetGene != g.Chromosome+"."+g.Name {
				res.Add(fmt.Sprintf("%s.references[%d]", genePath(gi), ri), fmt.Sprintf("missing dependency on %q for reference target", targetGene))
			}
			if _, ok := geneIndex[targetGene]; !ok {
				res.Add(fmt.Sprintf("%s.references[%d]", genePath(gi), ri), fmt.Sprintf("reference target gene %q not found in genome", targetGene))
			}
		}

		// entity fields of type reference
		for _, e := range g.Entities {
			for fname, fdef := range e.Fields {
				if fdef.Type == "reference" && fdef.Reference != "" {
					targetGene, err := loader.GenePartFromEntityRef(fdef.Reference)
					if err != nil {
						res.Add(fmt.Sprintf("%s.entities[%s].fields[%s]", genePath(gi), e.Name, fname), err.Error())
						continue
					}
					if !declaresDependency(g.Dependencies, targetGene) && targetGene != g.Chromosome+"."+g.Name {
						res.Add(fmt.Sprintf("%s.entities[%s].fields[%s]", genePath(gi), e.Name, fname), fmt.Sprintf("missing dependency on %q for field reference", targetGene))
					}
				}
			}
		}

		// capability inputs/outputs that are reference types
		for _, c := range g.Capabilities {
			for fname, fdef := range c.Inputs {
				checkCapRef(res, g, c.Name, "inputs", fname, fdef)
			}
			for fname, fdef := range c.Outputs {
				checkCapRef(res, g, c.Name, "outputs", fname, fdef)
			}
		}
	}
}

func checkCapRef(res *core.Result, g loader.ComposedGene, capName, ioKind, fname string, fdef loader.FieldDefinition) {
	if fdef.Type != "reference" || fdef.Reference == "" {
		return
	}
	targetGene, err := loader.GenePartFromEntityRef(fdef.Reference)
	if err != nil {
		res.Add(fmt.Sprintf("gene[%s.%s].capabilities[%s].%s[%s]", g.Chromosome, g.Name, capName, ioKind, fname), err.Error())
		return
	}
	if !declaresDependency(g.Dependencies, targetGene) && targetGene != g.Chromosome+"."+g.Name {
		res.Add(fmt.Sprintf("gene[%s.%s].capabilities[%s].%s[%s]", g.Chromosome, g.Name, capName, ioKind, fname),
			fmt.Sprintf("missing dependency on %q for capability reference", targetGene))
	}
}

func declaresDependency(deps []string, target string) bool {
	for _, d := range deps {
		if strings.EqualFold(d, target) {
			return true
		}
	}
	return false
}
