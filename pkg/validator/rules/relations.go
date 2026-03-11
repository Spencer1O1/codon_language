package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(relationsRules)
}

func relationsRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		codon, ok := gene.Codons["relations"].(map[string]any)
		if !ok {
			continue
		}
		for name, raw := range codon {
			rel, ok := raw.(map[string]any)
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_shape", Message: "relation must be an object", Gene: gene.Name, Codon: "relations"})
				continue
			}
			// ownership
			if ownRaw, ok := rel["ownership"]; ok {
				own, ok := ownRaw.(string)
				if !ok || (own != "from" && own != "to") {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "ownership_side_must_be_valid", Message: fmt.Sprintf("relation %s ownership must be 'from' or 'to'", name), Gene: gene.Name, Codon: "relations"})
				}
			}
			// cascade
			if casRaw, ok := rel["cascade"]; ok {
				cas, ok := casRaw.(string)
				if !ok || (cas != "cascade" && cas != "restrict" && cas != "nullify") {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "cascade_value_allowed", Message: fmt.Sprintf("relation %s cascade must be cascade|restrict|nullify", name), Gene: gene.Name, Codon: "relations"})
				}
			}
		}
	}
}
