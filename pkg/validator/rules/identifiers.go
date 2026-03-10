package rules

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkIdentifierPolicy)
}

func checkIdentifierPolicy(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		if err := loader.ValidateIdentifier("chromosome", g.Chromosome); err != nil {
			res.Add(genePath(gi), err.Error())
		}
		if err := loader.ValidateIdentifier("gene", g.Name); err != nil {
			res.Add(genePath(gi), err.Error())
		}
		for _, e := range g.Entities {
			if err := loader.ValidateIdentifier("entity", e.Name); err != nil {
				res.Add(genePath(gi)+".entities["+e.Name+"]", err.Error())
			}
			for fieldName := range e.Fields {
				if loader.IsReserved(fieldName) {
					res.Add(genePath(gi)+".entities["+e.Name+"].fields["+fieldName+"]", "field name is reserved word")
				}
			}
		}
		for _, c := range g.Capabilities {
			if err := loader.ValidateIdentifier("capability", c.Name); err != nil {
				res.Add(genePath(gi)+".capabilities["+c.Name+"]", err.Error())
			}
			for fname := range c.Inputs {
				if loader.IsReserved(fname) {
					res.Add(genePath(gi)+".capabilities["+c.Name+"].inputs["+fname+"]", "input name is reserved word")
				}
			}
			for fname := range c.Outputs {
				if loader.IsReserved(fname) {
					res.Add(genePath(gi)+".capabilities["+c.Name+"].outputs["+fname+"]", "output name is reserved word")
				}
			}
		}

		// traits/codons names are identifier-typed in source; here just enforce reserved words avoidance for traits names.
		for _, t := range g.Traits {
			if loader.IsReserved(t) {
				res.Add(genePath(gi)+".traits", "trait name is reserved word")
			}
		}
	}
}
