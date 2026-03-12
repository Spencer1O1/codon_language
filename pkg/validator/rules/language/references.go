package language

import (
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

var identRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

func init() {
	core.RegisterWithGroup("language", referenceRules)
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

// checkRef validates syntax and existence.
func checkRef(ref string, genome *loader.Genome, gene loader.Gene, codon string, res *core.Result) {
	parts := strings.Split(ref, ".")
	for _, p := range parts {
		if !identRe.MatchString(p) {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_must_match_pattern", Message: "ref must match identifier pattern", Gene: gene.Name, Codon: codon})
			return
		}
	}
	// quick over-qualification: if last segment matches a field in current gene
	if len(parts) > 1 {
		if ents, ok := gene.Codons["entities"].(map[string]any); ok {
			if _, ok := ents[parts[len(parts)-1]]; ok {
				res.Add(core.Issue{Severity: core.SeverityWarn, Code: "ref_overqualified", Message: "ref is over-qualified; use shortest form", Gene: gene.Name, Codon: codon})
				return
			}
		}
	}
	if resolveRef(ref, genome, gene) {
		return
	}
	if overQualifiedRef(ref, genome, gene) {
		res.Add(core.Issue{Severity: core.SeverityWarn, Code: "ref_overqualified", Message: "ref is over-qualified; use shortest form", Gene: gene.Name, Codon: codon})
		return
	}
	res.Add(core.Issue{Severity: core.SeverityError, Code: "ref_target_must_exist", Message: "ref target not found: " + ref, Gene: gene.Name, Codon: codon})
}

// resolveRef implements shortest path: name -> same codon; gene.name -> gene codon; chromosome.gene.name -> other chromosome.
func resolveRef(ref string, genome *loader.Genome, gene loader.Gene) bool {
	parts := strings.Split(ref, ".")
	switch len(parts) {
	case 1:
		// same codon name (entities capability fields)
		for _, g := range genome.Genes {
			if g.Name == gene.Name && g.Chromosome == gene.Chromosome {
				if _, ok := g.Codons["entities"].(map[string]any)[parts[0]]; ok {
					return true
				}
			}
		}
	case 2:
		gname := parts[0]
		field := parts[1]
		for _, g := range genome.Genes {
			if g.Name == gname {
				if m, ok := g.Codons["entities"].(map[string]any); ok {
					if _, ok := m[field]; ok {
						return true
					}
				}
			}
		}
	case 3:
		chrom, gname, field := parts[0], parts[1], parts[2]
		for _, g := range genome.Genes {
			if g.Chromosome == chrom && g.Name == gname {
				if m, ok := g.Codons["entities"].(map[string]any); ok {
					if _, ok := m[field]; ok {
						return true
					}
				}
			}
		}
	}
	return false
}

// overQualifiedRef returns true if a shorter equivalent reference exists.
func overQualifiedRef(ref string, genome *loader.Genome, gene loader.Gene) bool {
	parts := strings.Split(ref, ".")
	// chrom.gene.field -> gene.field only when chrom matches current chromosome
	if len(parts) == 3 && parts[0] == gene.Chromosome {
		if resolveRef(parts[1]+"."+parts[2], genome, gene) {
			return true
		}
	}
	// gene.field -> field only when gene matches current gene and chromosome
	if len(parts) == 2 && parts[0] == gene.Name {
		if resolveRef(parts[1], genome, gene) {
			return true
		}
	}
	return false
}
