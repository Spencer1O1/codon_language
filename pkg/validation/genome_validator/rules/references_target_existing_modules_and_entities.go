package rules

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type ReferencesTargetExistingModulesAndEntitiesRule struct{}

func (ReferencesTargetExistingModulesAndEntitiesRule) Name() string {
	return "references-target-existing-modules-and-entities"
}

func (ReferencesTargetExistingModulesAndEntitiesRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for moduleName, mod := range g.Modules {
		for i, ref := range mod.References {
			if _, ok := mod.Entities[ref.From]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "reference_missing_source_entity",
					Path:     fmt.Sprintf("module/%s/references/%d", moduleName, i),
					Message:  fmt.Sprintf("reference source entity %q does not exist in module %q", ref.From, moduleName),
				})
			}

			parts := strings.Split(ref.To, ".")
			if len(parts) != 2 {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "reference_invalid_target",
					Path:     fmt.Sprintf("module/%s/references/%d", moduleName, i),
					Message:  fmt.Sprintf("reference target %q must be in module.entity form", ref.To),
				})
				continue
			}

			targetModuleName := parts[0]
			targetEntityName := parts[1]

			targetModule, ok := g.Modules[targetModuleName]
			if !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "reference_missing_target_module",
					Path:     fmt.Sprintf("module/%s/references/%d", moduleName, i),
					Message:  fmt.Sprintf("reference target module %q does not exist", targetModuleName),
				})
				continue
			}

			if _, ok := targetModule.Entities[targetEntityName]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "reference_missing_target_entity",
					Path:     fmt.Sprintf("module/%s/references/%d", moduleName, i),
					Message:  fmt.Sprintf("reference target entity %q does not exist in module %q", targetEntityName, targetModuleName),
				})
			}
		}
	}

	return findings
}
