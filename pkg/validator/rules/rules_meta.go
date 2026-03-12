package rules

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.Register(rulesMeta) }

// rulesMeta enforces shape of rules codon and respects scope (defaults to validator).
func rulesMeta(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		rulesCodon, ok := gene.Codons["implementation"].(map[string]any)
		if !ok {
			continue
		}
		rulesRaw, ok := rulesCodon["rules"]
		if !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "rules_required", Message: "implementation.rules must be present", Gene: gene.Name, Codon: "implementation"})
			continue
		}
		list, ok := rulesRaw.([]any)
		if !ok || len(list) == 0 {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "rules_required", Message: "implementation.rules must be a non-empty list", Gene: gene.Name, Codon: "implementation"})
			continue
		}
		for _, r := range list {
			m, ok := r.(map[string]any)
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "code_nonempty", Message: "rule must be an object", Gene: gene.Name, Codon: "implementation"})
				continue
			}
			code, _ := m["code"].(string)
			if code == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "code_nonempty", Message: "rule code must be non-empty", Gene: gene.Name, Codon: "implementation"})
			}
			msg, _ := m["message"].(string)
			if msg == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "message_nonempty", Message: "rule message must be non-empty", Gene: gene.Name, Codon: "implementation"})
			}
			if sev, ok := m["severity"].(string); !ok || (sev != "error" && sev != "warn" && sev != "info") {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "severity_enum", Message: "rule severity must be error|warn|info", Gene: gene.Name, Codon: "implementation"})
			}
		}
	}
}
