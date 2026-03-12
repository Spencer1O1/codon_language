package language

import (
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

var identAllowed = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func init() {
	core.RegisterWithGroup("language", identifierRules)
}

func identifierRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	reserved := map[string]bool{
		"string": true, "number": true, "boolean": true, "uuid": true, "datetime": true, "json": true, "yaml": true, "ref": true, "TypeExpr": true, "Regex": true, "any": true, "object": true, "object_key": true, "object_value": true, "field": true, "primitive": true,
	}
	for _, gene := range g.Genes {
		if reserved[gene.Name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "gene name uses reserved identifier", Gene: gene.Name})
		}
		for codonName, codon := range gene.Codons {
			if m, ok := codon.(map[string]any); ok {
				for name := range m {
					if !identAllowed.MatchString(name) {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_syntax", Message: "identifiers must be lower_snake_case", Gene: gene.Name, Codon: codonName})
					}
					if reserved[name] {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_reserved", Message: "identifiers may not use reserved names (exported primitive nucleotypes or keywords)", Gene: gene.Name, Codon: codonName})
					}
					if strings.TrimSpace(name) == "" {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "identifier_nonempty", Message: "identifiers must be non-empty", Gene: gene.Name, Codon: codonName})
					}
				}
			}
		}
	}
}
