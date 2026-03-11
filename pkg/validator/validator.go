package validator

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
	_ "github.com/Spencer1O1/codon-language/pkg/validator/rules" // register rules via init
)

// Validate runs all registered rules against the loaded genome.
func Validate(g *loader.Genome, env map[string]nt.TypeNode) core.Result {
	res := core.Result{}
	for _, rule := range core.All() {
		rule(g, env, &res)
	}
	return res
}
