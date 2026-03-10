package validator

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
	"github.com/Spencer1O1/codon-language/pkg/validator/rules"
)

// Re-export types for convenience.
type (
	Result = core.Result
	Rule   = core.Rule
)

// Validate applies all registered rules to the composed genome.
func Validate(genome *loader.ComposedGenome) *Result {
	res := &core.Result{}
	for _, r := range core.Registry() {
		r(genome, res)
	}
	return res
}

// Touch the rules package so its init() registers all rules.
var _ = rules.All
