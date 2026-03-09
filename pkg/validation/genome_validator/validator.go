package genome_validator

import (
	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

/**
 * What to check:
 * 1. module dependencies exist
 * 2. relation entities exist
 * 3. reference targets exist
 * 4. no circular module dependencies
 * 5. entity names unique per module
 */

type Validator struct {
	rules []validation.Rule
}

func NewValidator(rules ...validation.Rule) *Validator {
	if len(rules) == 0 {
		rules = DefaultRules()
	}

	return &Validator{
		rules: rules,
	}
}

func (v *Validator) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	for _, rule := range v.rules {
		findings = append(findings, rule.Validate(g)...)
	}

	return findings
}
