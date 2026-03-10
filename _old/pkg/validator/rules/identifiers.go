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
			res.AddWithSeverity(core.Error, genePath(gi), err.Error())
		}
		if loader.IsReserved(g.Name) {
			res.AddWithSeverity(core.Warning, genePath(gi), "gene name is reserved word")
		}
		for _, e := range g.Entities {
			if err := loader.ValidateIdentifier("entity", e.Name); err != nil {
				res.AddWithSeverity(core.Error, genePath(gi)+".entities["+e.Name+"]", err.Error())
			}
			if loader.IsReserved(e.Name) {
				res.AddWithSeverity(core.Warning, genePath(gi)+".entities["+e.Name+"]", "entity name is reserved word")
			}
			for fieldName := range e.Fields {
				if loader.IsReserved(fieldName) {
					res.AddWithSeverity(core.Warning, genePath(gi)+".entities["+e.Name+"].fields["+fieldName+"]", "field name is reserved word")
				}
			}
		}
		for _, c := range g.Capabilities {
			if err := loader.ValidateIdentifier("capability", c.Name); err != nil {
				res.AddWithSeverity(core.Error, genePath(gi)+".capabilities["+c.Name+"]", err.Error())
			}
			if loader.IsReserved(c.Name) {
				res.AddWithSeverity(core.Warning, genePath(gi)+".capabilities["+c.Name+"]", "capability name is reserved word")
			}
			for fname := range c.Inputs {
				if loader.IsReserved(fname) {
					res.AddWithSeverity(core.Warning, genePath(gi)+".capabilities["+c.Name+"].inputs["+fname+"]", "input name is reserved word")
				}
			}
			for fname := range c.Outputs {
				if loader.IsReserved(fname) {
					res.AddWithSeverity(core.Warning, genePath(gi)+".capabilities["+c.Name+"].outputs["+fname+"]", "output name is reserved word")
				}
			}
		}

		for _, r := range g.Relations {
			if loader.IsReserved(r.Name) {
				res.AddWithSeverity(core.Warning, genePath(gi)+".relations["+r.Name+"]", "relation name is reserved word")
			}
		}

		for _, r := range g.References {
			if loader.IsReserved(r.Name) {
				res.AddWithSeverity(core.Warning, genePath(gi)+".references["+r.Name+"]", "reference name is reserved word")
			}
		}

		// traits/codons names are identifier-typed in source; here just enforce reserved words avoidance for traits names.
		for _, t := range g.Traits {
			if loader.IsReserved(t) {
				res.AddWithSeverity(core.Warning, genePath(gi)+".traits", "trait name is reserved word")
			}
		}
	}

	if loader.IsReserved(genome.Project.Name) {
		res.AddWithSeverity(core.Warning, "genome.project.name", "project name is reserved word")
	}
}
