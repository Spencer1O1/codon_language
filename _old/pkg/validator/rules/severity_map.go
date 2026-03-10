package rules

import "github.com/Spencer1O1/codon-language/pkg/validator/core"

var severityByCategory = map[core.RuleCategory]core.Severity{
	core.CategoryStructural:   core.Error,
	core.CategoryIdentifiers:  core.Error,
	core.CategoryDependencies: core.Error,
	core.CategoryReferences:   core.Error,
	core.CategoryRelations:    core.Error,
	core.CategoryUniqueness:   core.Error,
	core.CategoryTraits:       core.Warning,
	core.CategoryReserved:     core.Warning,
}

func severityFor(cat core.RuleCategory, fallback core.Severity) core.Severity {
	if sev, ok := severityByCategory[cat]; ok {
		return sev
	}
	return fallback
}
