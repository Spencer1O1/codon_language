package language

import (
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("language", capabilityRules) }

func capabilityRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		codon, ok := gene.Codons["capabilities"].(map[string]any)
		if !ok {
			continue
		}
		seen := map[string]bool{}
		for name, raw := range codon {
			if strings.TrimSpace(name) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "capability_key_required", Message: "capability names must be non-empty", Gene: gene.Name, Codon: "capabilities"})
				continue
			}
			if seen[name] {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "capability_name_unique_within_capabilities", Message: "capability names must be unique within the capabilities codon", Gene: gene.Name, Codon: "capabilities"})
				continue
			}
			seen[name] = true
			cap, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if effects, ok := cap["effects"].([]any); !ok || len(effects) == 0 {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "effects_required", Message: "capability.effects must be a non-empty list", Gene: gene.Name, Codon: "capabilities"})
			}
		}
	}
}
