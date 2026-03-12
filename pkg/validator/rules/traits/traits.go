package traits

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() {
	core.RegisterWithGroup("traits", traitRules)
}

func traitRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	// gene traits
	for _, gene := range g.Genes {
		if traitsRaw, ok := gene.Codons["traits"]; ok {
			list, ok := traitsRaw.([]any)
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "trait_is_string", Message: "traits must be a list of strings", Gene: gene.Name, Codon: "traits"})
				continue
			}
			seen := map[string]map[string]any{}
			for _, tr := range list {
				name, ok := tr.(string)
				if !ok || strings.TrimSpace(name) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "trait_is_string", Message: "trait entries must be strings", Gene: gene.Name, Codon: "traits"})
					continue
				}
				pattern := filepath.Join(g.Root, "traits", "gene", name, "trait.yaml")
				if _, err := os.Stat(pattern); err != nil {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "gene_trait_file_exists", Message: "gene trait file not found: " + pattern, Gene: gene.Name, Codon: "traits"})
					continue
				}
				if _, exists := seen[name]; exists {
					res.Add(core.Issue{Severity: core.SeverityWarn, Code: "trait_conflict_resolution", Message: fmt.Sprintf("duplicate trait %s on gene %s; authored definitions win", name, gene.Name), Gene: gene.Name, Codon: "traits"})
				}
				seen[name] = nil
			}
		}
	}

	// chromosome traits
	for chrom, gene := range groupByChromosome(g.Genes) {
		if ctraits, ok := g.Manifest["chromosome_traits"].(map[string]any); ok {
			if nameRaw, ok := ctraits[chrom]; ok {
				name, _ := nameRaw.(string)
				if name != "" {
					pattern := filepath.Join(g.Root, "traits", "chromosome", name, "trait.yaml")
					if _, err := os.Stat(pattern); err != nil {
						res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_trait_file_exists", Message: "chromosome trait file not found: " + pattern, Gene: gene, Codon: "traits"})
					}
				}
			}
		}
	}
}

func groupByChromosome(genes []loader.Gene) map[string]string {
	m := map[string]string{}
	for _, g := range genes {
		m[g.Chromosome] = g.Name
	}
	return m
}
