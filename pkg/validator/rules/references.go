package rules

import (
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

var identRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

func init() {
	core.Register(referenceRules)
}

func referenceRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		for codonName, val := range gene.Codons {
			walkRefs(val, g, gene, codonName, res)
		}
	}
}

func walkRefs(v any, genome *loader.Genome, gene loader.Gene, codon string, res *core.Result) {
	switch t := v.(type) {
	case map[string]any:
		if refVal, ok := t["ref"]; ok {
			if refStr, ok := refVal.(string); ok {
				checkRef(refStr, genome, gene, codon, res)
			} else {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_must_match_pattern", Message: "ref must be a string", Gene: gene.Name, Codon: codon})
			}
		}
		for _, vv := range t {
			walkRefs(vv, genome, gene, codon, res)
		}
	case []any:
		for _, vv := range t {
			walkRefs(vv, genome, gene, codon, res)
		}
	}
}

func checkRef(ref string, genome *loader.Genome, gene loader.Gene, codon string, res *core.Result) {
	parts := strings.Split(ref, ".")
	if len(parts) < 1 || len(parts) > 4 {
		res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_must_match_pattern", Message: "ref must have 1-4 dot-separated identifiers", Gene: gene.Name, Codon: codon})
		return
	}
	for _, p := range parts {
		if !identRe.MatchString(p) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_must_match_pattern", Message: "ref identifiers must match [A-Za-z][A-Za-z0-9_]*", Gene: gene.Name, Codon: codon})
			return
		}
	}

	// attempt resolution when fully qualified
	if len(parts) == 4 {
		chrom, geneName, codonName, entry := parts[0], parts[1], parts[2], parts[3]
		foundGene := false
		foundCodon := false
		foundEntry := false
		for _, g := range genome.Genes {
			if g.Chromosome == chrom && g.Name == geneName {
				foundGene = true
				if c, ok := g.Codons[codonName].(map[string]any); ok {
					foundCodon = true
					if _, ok := c[entry]; ok {
						foundEntry = true
					}
				}
			}
		}
		if !foundGene {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref gene not found: " + ref, Gene: gene.Name, Codon: codon})
			return
		}
		if !foundCodon {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref codon not found: " + ref, Gene: gene.Name, Codon: codon})
			return
		}
		if !foundEntry {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref entry not found: " + ref, Gene: gene.Name, Codon: codon})
			return
		}
	}
}
