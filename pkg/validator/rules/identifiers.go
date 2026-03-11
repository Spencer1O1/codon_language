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

func identifierRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	// chromosome and gene names from paths
	for _, gene := range g.Genes {
		if !identPattern.MatchString(gene.Chromosome) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "chromosome name must be lower_snake", Gene: gene.Name, Codon: ""})
		}
		if !identPattern.MatchString(strings.ToLower(gene.Name)) { // allow dash? spec uses lower_snake
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "gene name must be lower_snake", Gene: gene.Name, Codon: ""})
		}
		for codonName, val := range gene.Codons {
			if !identPattern.MatchString(codonName) {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "codon name must be lower_snake", Gene: gene.Name, Codon: codonName})
			}
			if m, ok := val.(map[string]any); ok {
				for entry := range m {
					if !identPattern.MatchString(entry) {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_invalid", Message: "entry name must be lower_snake", Gene: gene.Name, Codon: codonName})
					}
				}
			}
		}
	}
}
