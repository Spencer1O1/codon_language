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
	boundCaps := map[string]string{} // capability -> projection name
	allCaps := collectCapabilities(g)
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
		tgtVal, tgtOK := m["target"]
		tgt, _ := tgtVal.(string)
		if !tgtOK {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_target_exists", Message: "projection target must exist in targets.yaml", Codon: "projections"})
		} else if _, ok := tgtVal.(string); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_target_string", Message: "projection.target must be a string", Codon: "projections"})
		} else if tgt == "" || !targets[tgt] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_target_exists", Message: "projection target must exist in targets.yaml", Codon: "projections"})
		}
		bindingVal, bindingOK := m["binding"]
		if !bindingOK {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_binding_named", Message: "projection binding is required", Codon: "projections"})
		} else {
			if _, ok := bindingVal.(string); !ok || bindingVal.(string) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_binding_string", Message: "projection.binding must be a string", Codon: "projections"})
			}
		}
		selAny := false
		for _, key := range []string{"capabilities", "entities", "relations"} {
			if v, ok := m[key]; ok {
				switch vv := v.(type) {
				case string:
					if key == "capabilities" && vv == "*" {
						selAny = true
						for capName := range allCaps {
							if prev, ok := boundCaps[capName]; ok {
								res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_capability_unique", Message: "capability bound to multiple projections", Codon: "projections"})
								_ = prev
							}
							boundCaps[capName] = name
						}
					} else {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_selectors_type", Message: "selectors must be lists of strings or \"*\" for capabilities", Codon: "projections"})
					}
				case []any:
					if len(vv) > 0 {
						selAny = true
					}
					for _, entry := range vv {
						if _, ok := entry.(string); !ok {
							res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_selector_elements_string", Message: "selector entries must be strings", Codon: "projections"})
						}
						if key == "capabilities" {
							for _, entry := range vv {
								if capName, ok := entry.(string); ok {
									if prev, ok := boundCaps[capName]; ok {
										res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_capability_unique", Message: "capability bound to multiple projections", Codon: "projections"})
										_ = prev
									}
									boundCaps[capName] = name
								}
							}
						}
					}
				default:
					res.Add(core.Issue{Severity: core.SeverityError, Code: "projection_selectors_type", Message: "selectors must be lists of strings or \"*\" for capabilities", Codon: "projections"})
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

	// coverage: info-level for now
	for capName := range allCaps {
		if _, ok := boundCaps[capName]; !ok {
			res.Add(core.Issue{Severity: core.SeverityInfo, Code: "projection_capability_coverage", Message: "capability is not bound to any projection", Codon: "projections"})
		}
	}
}

// collectCapabilities gathers all capability names in the genome.
func collectCapabilities(g *loader.Genome) map[string]struct{} {
	caps := map[string]struct{}{}
	for _, gene := range g.Genes {
		if m, ok := gene.Codons["capabilities"].(map[string]any); ok {
			for name := range m {
				caps[name] = struct{}{}
			}
		}
	}
	return caps
}
