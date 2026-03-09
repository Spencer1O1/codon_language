package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type RelationsReferenceExistingEntitiesRule struct{}

func (RelationsReferenceExistingEntitiesRule) Name() string {
	return "relations-reference-existing-entities"
}

func (RelationsReferenceExistingEntitiesRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for moduleName, mod := range g.Modules {
		for i, rel := range mod.Relations {
			if _, ok := mod.Entities[rel.From]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "relation_missing_source_entity",
					Path:     fmt.Sprintf("module/%s/relations/%d", moduleName, i),
					Message:  fmt.Sprintf("relation source entity %q does not exist in module %q", rel.From, moduleName),
				})
			}

			if _, ok := mod.Entities[rel.To]; !ok {
				findings = append(findings, validation.Finding{
					Severity: validation.SeverityError,
					Code:     "relation_missing_target_entity",
					Path:     fmt.Sprintf("module/%s/relations/%d", moduleName, i),
					Message:  fmt.Sprintf("relation target entity %q does not exist in module %q", rel.To, moduleName),
				})
			}
		}
	}

	return findings
}
