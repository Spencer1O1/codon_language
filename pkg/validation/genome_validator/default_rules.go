package genome_validator

import (
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
	"github.com/Spencer1O1/codon-language/internal/domain/validation/genome_validator/rules"
)

func DefaultRules() []validation.Rule {
	return []validation.Rule{
		rules.DependenciesExistRule{},
		rules.NoSelfDependencyRule{},
		rules.NoCircularDependenciesRule{},
		rules.RelationsReferenceExistingEntitiesRule{},
		rules.ReferencesTargetExistingModulesAndEntitiesRule{},
		rules.ReferencesRequireDeclaredDependencyRule{},
	}
}
