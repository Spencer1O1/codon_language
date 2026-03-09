package rules

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type ReferencesRequireDeclaredDependencyRule struct{}

func (ReferencesRequireDeclaredDependencyRule) Name() string {
	return "references-require-declared-dependency"
}

func (ReferencesRequireDeclaredDependencyRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for moduleName, mod := range g.Modules {
		deps := make(map[string]struct{}, len(mod.Dependencies))
		for _, dep := range mod.Dependencies {
			deps[dep] = struct{}{}
		}

		for i, ref := range mod.References {
			parts := strings.Split(ref.To, ".")
			if len(parts) != 2 {
				// Let the "invalid target format" rule handle this
				continue
			}

			targetModuleName := parts[0]

			// Same-module references do not require a declared dependency
			if targetModuleName == moduleName {
				continue
			}

			if _, ok := deps[targetModuleName]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "reference_missing_declared_dependency",
					Path:     fmt.Sprintf("module/%s/references/%d", moduleName, i),
					Message: fmt.Sprintf(
						"module %q references %q but does not declare dependency on module %q",
						moduleName,
						ref.To,
						targetModuleName,
					),
				})
			}
		}
	}

	return findings
}
