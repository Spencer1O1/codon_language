package expression

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("expression", stylesRules) }

func stylesRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil || g.Expression.Styles == nil {
		return
	}
	smap := g.Expression.Styles
	if smap == nil {
		return
	}
	seen := map[string]bool{}
	for name, raw := range smap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "style_names_unique", Message: "style names must be unique", Codon: "styles"})
			continue
		}
		if _, ok := raw.(map[string]any); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "styles_shape_map", Message: "styles.yaml must be a map of style_name to object", Codon: "styles"})
			continue
		}
		seen[name] = true
		if m, ok := raw.(map[string]any); ok {
			if v, vok := m["version"]; vok {
				if s, ok := v.(string); ok && s == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "style_version_recommended", Message: "style version should be non-empty when provided", Codon: "styles"})
				}
			}
		}
	}
}
