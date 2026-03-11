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
		if gene.Chromosome == "language" {
			continue // skip doc/examples genes
		}
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
	switch len(parts) {
	case 1:
		if !hasEntry(&gene, codon, parts[0]) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref entry not found: " + ref, Gene: gene.Name, Codon: codon})
		}
	case 2:
		if !hasEntry(&gene, parts[0], parts[1]) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref entry not found: " + ref, Gene: gene.Name, Codon: codon})
		} else if hasEntry(&gene, codon, parts[1]) {
			res.Add(core.Issue{Severity: core.SeverityWarn, Code: "ref_overqualified", Message: "reference could be shortened to " + parts[1], Gene: gene.Name, Codon: codon})
		}
	case 3:
		target := findGene(genome, gene.Chromosome, parts[0])
		if !hasEntry(target, parts[1], parts[2]) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref entry not found: " + ref, Gene: gene.Name, Codon: codon})
		} else if hasEntry(&gene, parts[1], parts[2]) {
			res.Add(core.Issue{Severity: core.SeverityWarn, Code: "ref_overqualified", Message: "reference could be shortened to " + parts[1] + "." + parts[2], Gene: gene.Name, Codon: codon})
		}
	case 4:
		target := findGene(genome, parts[0], parts[1])
		if !hasEntry(target, parts[2], parts[3]) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref entry not found: " + ref, Gene: gene.Name, Codon: codon})
		} else if gene.Chromosome == parts[0] && gene.Name == parts[1] && hasEntry(&gene, parts[2], parts[3]) {
			res.Add(core.Issue{Severity: core.SeverityWarn, Code: "ref_overqualified", Message: "reference could be shortened to " + parts[2] + "." + parts[3], Gene: gene.Name, Codon: codon})
		}
	}
}

func findGene(genome *loader.Genome, chrom, geneName string) *loader.Gene {
	for i, g := range genome.Genes {
		if g.Chromosome == chrom && g.Name == geneName {
			return &genome.Genes[i]
		}
	}
	return nil
}

func hasEntry(g *loader.Gene, codonName, entry string) bool {
	if g == nil {
		return false
	}
	c, ok := g.Codons[codonName].(map[string]any)
	if !ok {
		return false
	}
	_, ok = c[entry]
	return ok
}
