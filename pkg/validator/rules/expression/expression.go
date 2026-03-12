package expression

import (
	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() {
	core.RegisterWithGroup("expression", validateExpressionAssets)
}

func validateExpressionAssets(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Expression == nil {
		return // expression is optional; use --require-expression policy elsewhere if needed
	}
	checkMap := func(name string, m map[string]any) {
		// already map from loader; nothing deeper yet
		if m == nil {
			return
		}
	}
	checkMap("targets", g.Expression.Targets)
	checkMap("projections", g.Expression.Projections)
	checkMap("styles", g.Expression.Styles)
	checkMap("templates", g.Expression.Templates)
}
