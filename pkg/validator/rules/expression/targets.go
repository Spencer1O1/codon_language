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
	tmap := g.Expression.Targets
	if tmap == nil {
		return
	}
	seen := map[string]bool{}
	for name, raw := range tmap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "target_names_unique", Message: "target names must be unique", Codon: "targets"})
			continue
		}
		if _, ok := raw.(map[string]any); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "targets_shape_map", Message: "targets.yaml must be a map of target_name to object", Codon: "targets"})
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
		if kindVal, ok := m["kind"]; ok {
			if _, ok := kindVal.(string); !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "target_kind_must_be_string", Message: "target.kind must be a string", Codon: "targets"})
			}
		}
		if stack != "" {
			if _, ok := m["stack"].(string); !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "target_stack_must_be_string", Message: "target.stack must be a string", Codon: "targets"})
			}
		}
		if out, ok := m["output_root"]; ok {
			if _, ok := out.(string); !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "target_output_root_string", Message: "output_root must be a string when provided", Codon: "targets"})
			}
		}
		if ow, ok := m["overwrite"]; ok {
			if ovs, ok := ow.(string); ok && ovs != "" {
				switch ovs {
				case "safe", "force":
				default:
					res.Add(core.Issue{Severity: core.SeverityError, Code: "target_overwrite_value", Message: "overwrite must be safe|force when provided", Codon: "targets"})
				}
			} else {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "target_overwrite_value", Message: "overwrite must be a string (safe|force) when provided", Codon: "targets"})
			}
		}
		if _, ok := m["output_root"]; !ok {
			res.Add(core.Issue{Severity: core.SeverityInfo, Code: "target_output_root_recommended", Message: "output_root is recommended for targets", Codon: "targets"})
		}
	}
}
