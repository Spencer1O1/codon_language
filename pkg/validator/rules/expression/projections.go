package expression

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("expression", projectionsRules) }

func projectionsRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil || g.Expression.Projections == nil {
		return
	}
	pmap := g.Expression.Projections
	if pmap == nil {
		return
	}
	targets := map[string]bool{}
	if g.Expression.Targets != nil {
		for k := range g.Expression.Targets {
			targets[k] = true
		}
	}
	seen := map[string]bool{}
	for name, raw := range pmap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_names_unique", Message: "projection names must be unique", Codon: "projections"})
			continue
		}
		if _, ok := raw.(map[string]any); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projections_shape_map", Message: "projections.yaml must be a map of projection_name to object", Codon: "projections"})
			continue
		}
		seen[name] = true
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		tgt, _ := m["target"].(string)
		if tgt == "" || !targets[tgt] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_target_exists", Message: "projection target must exist in targets.yaml", Codon: "projections"})
		}
		if _, ok := m["binding"]; !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_binding_named", Message: "projection binding is required", Codon: "projections"})
		}
		selAny := false
		for _, key := range []string{"capabilities", "entities", "relations"} {
			if v, ok := m[key]; ok {
				if list, ok := v.([]any); ok && len(list) > 0 {
					selAny = true
				}
			}
		}
		if !selAny {
			if v, ok := m["capabilities"].(string); ok && v == "*" {
				selAny = true
			}
		}
		if !selAny {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_selectors_nonempty", Message: "projection must select capabilities/entities/relations or '*'", Codon: "projections"})
		}
	}
}
