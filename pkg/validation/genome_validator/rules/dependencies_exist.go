package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type DependenciesExistRule struct{}

func (DependenciesExistRule) Name() string { return "dependencies-exist" }

func (DependenciesExistRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for moduleName, mod := range g.Modules {
		for _, dep := range mod.Dependencies {
			if _, ok := g.Modules[dep]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "missing_dependency",
					Path:     fmt.Sprintf("module/%s/dependencies", moduleName),
					Message:  fmt.Sprintf("module %q depends on unknown module %q", moduleName, dep),
				})
			}
		}
	}

	return findings
}
