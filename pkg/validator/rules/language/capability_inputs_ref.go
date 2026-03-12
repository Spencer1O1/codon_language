package language

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// ref_type_usage: if a field map has ref, its type/type_expr must be ref or absent.
func init() { core.RegisterWithGroup("language", refTypeUsageRules) }

func refTypeUsageRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for codonName, val := range gene.Codons {
			walkRefType(val, gene, codonName, res)
		}
	}
}

func walkRefType(v any, gene loader.Gene, codon string, res *core.Result) {
	switch t := v.(type) {
	case map[string]any:
		if _, ok := t["ref"]; ok {
			if te, ok := t["type_expr"].(string); ok && te != "ref" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_type_usage", Message: "fields using ref must declare type_expr: ref or omit type/type_expr", Gene: gene.Name, Codon: codon})
			}
			if tt, ok := t["type"].(string); ok && tt != "ref" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_type_usage", Message: "fields using ref must declare type: ref or omit type/type_expr", Gene: gene.Name, Codon: codon})
			}
		}
		for _, vv := range t {
			walkRefType(vv, gene, codon, res)
		}
	case []any:
		for _, vv := range t {
			walkRefType(vv, gene, codon, res)
		}
	}
}
