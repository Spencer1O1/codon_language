package expression

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("expression", templatesRules) }

func templatesRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil || g.Expression.Templates == nil {
		return
	}
	tmap, ok := g.Expression.Templates["templates"].(map[string]any)
	if !ok {
		return
	}
	seen := map[string]bool{}
	for name, raw := range tmap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_names_unique", Message: "template names must be unique", Codon: "templates"})
			continue
		}
		seen[name] = true
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		src, _ := m["source"].(string)
		if src == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_source_required", Message: "template source is required", Codon: "templates"})
		}
		if _, ok := m["checksum"]; !ok {
			res.Add(core.Issue{Severity: core.SeverityInfo, Code: "template_checksum_recommended", Message: "template checksum is recommended", Codon: "templates"})
		}
	}
}
