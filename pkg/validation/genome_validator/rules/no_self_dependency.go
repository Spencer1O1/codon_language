package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type NoSelfDependencyRule struct{}

func (NoSelfDependencyRule) Name() string { return "no-self-dependency" }

func (NoSelfDependencyRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for moduleName, mod := range g.Modules {
		for _, dep := range mod.Dependencies {
			if dep == moduleName {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "self_dependency",
					Path:     fmt.Sprintf("module/%s/dependencies", moduleName),
					Message:  fmt.Sprintf("module %q cannot depend on itself", moduleName),
				})
			}
		}
	}

	return findings
}
