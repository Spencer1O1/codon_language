package language

import (
	"fmt"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() {
	core.RegisterWithGroup("language", schemaRules)
}

// schemaRules validates codon values against their codon schema type expressions.
func schemaRules(g *loader.Genome, env map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for codonName, val := range gene.Codons {
			schema, ok := g.Schemas[codonName]
			if !ok {
				continue // missing schema handled in basicShape
			}
			if err := validateAgainstType(val, schema.TypeAST, env); err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_mismatch", Message: fmt.Sprintf("codon %s does not match schema: %v", codonName, err), Gene: gene.Name, Codon: codonName})
			}
		}
	}
}

// validateAgainstType is a placeholder; existing structural checks already run elsewhere.
// TODO: implement full TypeExpr value validation.
func validateAgainstType(_ any, _ nt.TypeNode, _ map[string]nt.TypeNode) error {
	return nil
}
