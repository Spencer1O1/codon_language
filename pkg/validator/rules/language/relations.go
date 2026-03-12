package language

import (
	"strings"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() {
	core.RegisterWithGroup("language", relationsRules)
}

func relationsRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		codon, ok := gene.Codons["relations"].(map[string]any)
		if !ok {
			continue
		}
		seen := map[string]bool{}
		for name, raw := range codon {
			if strings.TrimSpace(name) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_key_required", Message: "relation names must be non-empty", Gene: gene.Name, Codon: "relations"})
				continue
			}
			if seen[name] {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_name_unique_within_relations", Message: "relation names must be unique within the relations codon", Gene: gene.Name, Codon: "relations"})
				continue
			}
			seen[name] = true
			rmap, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			from, _ := rmap["from"].(string)
			to, _ := rmap["to"].(string)
			if strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_ref_required", Message: "relations must specify from and to", Gene: gene.Name, Codon: "relations"})
				continue
			}
			checkRelationEnd := func(ref string) {
				if overQualifiedRef(ref, g, gene) {
					res.Add(core.Issue{Severity: core.SeverityInfo, Code: "relation_overqualified", Message: "relation endpoint is over-qualified; use shortest form", Gene: gene.Name, Codon: "relations"})
					return
				}
				if resolveRef(ref, g, gene) {
					return
				}
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_target_must_exist", Message: "relation endpoints must exist", Gene: gene.Name, Codon: "relations"})
			}
			checkRelationEnd(from)
			checkRelationEnd(to)
			if ow, ok := rmap["ownership"].(string); ok && ow != "" && ow != "from" && ow != "to" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "ownership_side_must_be_valid", Message: "ownership must be from|to", Gene: gene.Name, Codon: "relations"})
			}
			if cs, ok := rmap["cascade"].(string); ok && cs != "" {
				switch cs {
				case "cascade", "restrict", "nullify":
				default:
					res.Add(core.Issue{Severity: core.SeverityError, Code: "cascade_value_allowed", Message: "cascade must be cascade|restrict|nullify", Gene: gene.Name, Codon: "relations"})
				}
			}
		}
	}
}
