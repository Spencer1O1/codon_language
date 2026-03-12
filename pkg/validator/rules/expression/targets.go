package expression

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("expression", targetsRules) }

func targetsRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil || g.Expression.Targets == nil {
		return
	}
	tmap, ok := g.Expression.Targets["targets"].(map[string]any)
	if !ok {
		return
	}
	seen := map[string]bool{}
	for name, raw := range tmap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "target_names_unique", Message: "target names must be unique", Codon: "targets"})
			continue
		}
		seen[name] = true
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		kind, _ := m["kind"].(string)
		stack, _ := m["stack"].(string)
		if kind == "" || stack == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "target_requires_kind_and_stack", Message: "target.kind and target.stack are required", Codon: "targets"})
		}
		if _, ok := m["output_root"]; !ok {
			res.Add(core.Issue{Severity: core.SeverityInfo, Code: "target_output_root_recommended", Message: "output_root is recommended for targets", Codon: "targets"})
		}
	}
}
