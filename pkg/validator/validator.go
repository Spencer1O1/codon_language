package validator

import (
	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
	_ "github.com/Spencer1O1/codon_language/pkg/validator/rules/expression" // placeholder group (future)
	_ "github.com/Spencer1O1/codon_language/pkg/validator/rules/language"   // register language rules
	_ "github.com/Spencer1O1/codon_language/pkg/validator/rules/manifest"   // register manifest rules
	_ "github.com/Spencer1O1/codon_language/pkg/validator/rules/traits"     // register trait rules
)

// Validate runs all registered rules against the loaded genome.
func Validate(g *loader.Genome, env map[string]nt.TypeNode) core.Result {
	res := core.Result{}
	// keep env in sync with genome TypeEnv (may be extended by traits)
	if g.TypeEnv != nil {
		env = g.TypeEnv
	}
	applyTraitInjection(g, &res)
	groups := []string{"manifest", "language", "traits", "expression"}
	for _, grp := range groups {
		for _, rule := range core.AllByGroup(grp) {
			rule(g, env, &res)
		}
	}
	// fallback: run any rules registered to other groups (if any)
	for _, grp := range core.Groups() {
		seen := false
		for _, fixed := range groups {
			if grp == fixed {
				seen = true
				break
			}
		}
		if seen {
			continue
		}
		for _, rule := range core.AllByGroup(grp) {
			rule(g, env, &res)
		}
	}
	return res
}
