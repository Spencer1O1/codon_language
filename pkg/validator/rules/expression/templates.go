package expression

import (
	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("expression", templatesRules) }

func templatesRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil || g.Expression.Templates == nil {
		return
	}
	tmap := g.Expression.Templates
	if tmap == nil {
		return
	}
	seen := map[string]bool{}
	for name, raw := range tmap {
		if seen[name] {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_names_unique", Message: "template names must be unique", Codon: "templates"})
			continue
		}
		if _, ok := raw.(map[string]any); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "templates_shape_map", Message: "templates.yaml must be a map of template_name to object", Codon: "templates"})
			continue
		}
		seen[name] = true
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		srcVal, srcOK := m["source"]
		src, _ := srcVal.(string)
		if !srcOK {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_source_required", Message: "template source is required", Codon: "templates"})
		} else if _, ok := srcVal.(string); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_source_string", Message: "template.source must be a string", Codon: "templates"})
		} else if src == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_source_required", Message: "template source is required", Codon: "templates"})
		}
		if _, ok := m["checksum"]; !ok {
			res.Add(core.Issue{Severity: core.SeverityInfo, Code: "template_checksum_recommended", Message: "template checksum is recommended", Codon: "templates"})
		} else if _, ok := m["checksum"].(string); !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "template_checksum_string", Message: "template.checksum must be a string", Codon: "templates"})
		}
		if v, ok := m["variables"]; ok {
			if _, ok := v.(map[string]any); !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "template_variables_map", Message: "template.variables must be a map", Codon: "templates"})
			}
		}
		if v, ok := m["postprocess"]; ok {
			list, ok := v.([]any)
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "template_postprocess_strings", Message: "template.postprocess must be a list of strings", Codon: "templates"})
			} else {
				for _, e := range list {
					if _, ok := e.(string); !ok {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "template_postprocess_strings", Message: "template.postprocess must be a list of strings", Codon: "templates"})
					}
				}
			}
		}
	}
}
