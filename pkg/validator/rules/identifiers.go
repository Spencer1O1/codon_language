package rules

import (
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

var identPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func init() {
	core.Register(identifierRules)
}

func identifierRules(g *loader.Genome, env map[string]nt.TypeNode, res *core.Result) {
	// chromosome and gene names from paths
	reserved := buildReserved()

	for _, gene := range g.Genes {
		if !identPattern.MatchString(gene.Chromosome) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "chromosome name must be lower_snake", Gene: gene.Name, Codon: ""})
		}
		if reserved[strings.ToLower(gene.Chromosome)] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "chromosome name is reserved", Gene: gene.Name, Codon: ""})
		}
		if !identPattern.MatchString(strings.ToLower(gene.Name)) { // allow dash? spec uses lower_snake
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "gene name must be lower_snake", Gene: gene.Name, Codon: ""})
		}
		if reserved[strings.ToLower(gene.Name)] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "gene name is reserved", Gene: gene.Name, Codon: ""})
		}
		for codonName, val := range gene.Codons {
			if !identPattern.MatchString(codonName) {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "codon name must be lower_snake", Gene: gene.Name, Codon: codonName})
			}
			if reserved[strings.ToLower(codonName)] {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "codon name is reserved", Gene: gene.Name, Codon: codonName})
			}
			if m, ok := val.(map[string]any); ok {
				for entry := range m {
					if !identPattern.MatchString(entry) {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "entry name must be lower_snake", Gene: gene.Name, Codon: codonName})
					}
					if reserved[strings.ToLower(entry)] {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "entry name is reserved", Gene: gene.Name, Codon: codonName})
					}
				}
			}
		}
	}
}

func buildReserved() map[string]bool {
	r := map[string]bool{}
	// exported primitive nucleotypes (hard-coded list from primitives.nucleotype)
	for _, name := range []string{"string", "number", "boolean", "uuid", "datetime", "json", "yaml", "ref", "regex", "typeexpr", "any", "object"} {
		r[name] = true
	}
	return r
}
