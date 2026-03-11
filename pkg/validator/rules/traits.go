package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(traitRules)
}

func traitRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	// genome traits in manifest
	if traitsRaw, ok := g.Manifest["traits"]; ok {
		if traits, ok := traitsRaw.(map[string]any); ok {
			for chrom, val := range traits {
				name, ok := val.(string)
				if !ok || strings.TrimSpace(name) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "traits_are_map", Message: "traits values must be trait names"})
					continue
				}
				if strings.TrimSpace(chrom) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_requires_chromosome", Message: "manifest traits must specify chromosome key"})
				}
				pattern := filepath.Join(g.Root, "traits", "genome", name+".yaml")
				matches, _ := filepath.Glob(pattern)
				if len(matches) == 0 {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_file_exists", Message: "manifest trait file not found: " + pattern})
				}
			}
		}
	}

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
				pattern := filepath.Join(g.Root, "traits", "gene", name+".yaml")
				matches, _ := filepath.Glob(pattern)
				if len(matches) == 0 {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "gene_trait_file_exists", Message: "gene trait file not found: " + pattern, Gene: gene.Name, Codon: "traits"})
					continue
				}
				// simplistic conflict detection: if multiple traits inject same codon key, require identical
				// (full trait injection not implemented; placeholder for policy)
				if _, exists := seen[name]; exists {
					res.Add(core.Issue{Severity: core.SeverityWarn, Code: "trait_conflict_authored_wins", Message: fmt.Sprintf("duplicate trait %s on gene %s", name, gene.Name), Gene: gene.Name, Codon: "traits"})
				} else {
					seen[name] = nil
				}
			}
		}
	}
}
