package rules

import (
    "path/filepath"
    "strings"
    "os"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(manifestRules)
}

// manifestRules enforces manifest-level rules documented in genome_manifest validation codon.
func manifestRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	if g.Manifest == nil {
		res.Add(core.Issue{Severity: core.SeverityError, Code: "manifest_file_exists", Message: "manifest (genome.yaml) missing"})
		return
	}

	// schema_version
	if sv, ok := g.Manifest["schema_version"].(string); !ok || strings.TrimSpace(sv) == "" {
		res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_version_required", Message: "schema_version is required"})
	}

	// project
	proj, ok := g.Manifest["project"].(map[string]any)
	if !ok {
		res.Add(core.Issue{Severity: core.SeverityError, Code: "project_name_required", Message: "project block is required"})
	} else {
		if name, ok := proj["name"].(string); !ok || strings.TrimSpace(name) == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "project_name_required", Message: "project.name is required"})
		}
		if typ, ok := proj["type"].(string); !ok || strings.TrimSpace(typ) == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "project_type_required", Message: "project.type is required"})
		}
	}

	// manifest filename check (best-effort path from loader caller)
	// Not robust without path info; assume loader root contains genome.yaml.
	// traits map: chromosome -> trait name (legacy global scope)
	if traitsRaw, ok := g.Manifest["traits"]; ok {
		traits, ok := traitsRaw.(map[string]any)
		if !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "traits_are_map", Message: "traits must be a map of chromosome_name to trait_name"})
		} else {
			for chrom, val := range traits {
				name, ok := val.(string)
				if !ok || strings.TrimSpace(name) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "traits_are_map", Message: "traits values must be trait names"})
					continue
				}
				if strings.TrimSpace(chrom) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_requires_chromosome", Message: "manifest traits must specify chromosome key"})
				}
				pattern := filepath.Join(g.Root, "traits", "genome", name, "trait.yaml")
				if _, err := os.Stat(pattern); err != nil {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_file_exists", Message: "manifest trait file not found: " + pattern})
				}
			}
		}
	}

	// chromosome_traits map: chromosome -> trait name
	if cTraitsRaw, ok := g.Manifest["chromosome_traits"]; ok {
		traits, ok := cTraitsRaw.(map[string]any)
		if !ok {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_traits_are_map", Message: "chromosome_traits must be a map of chromosome_name to trait_name"})
		} else {
			for chrom, val := range traits {
				name, ok := val.(string)
				if !ok || strings.TrimSpace(name) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_traits_are_map", Message: "chromosome_traits values must be trait names"})
					continue
				}
				if strings.TrimSpace(chrom) == "" {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_trait_requires_chromosome", Message: "chromosome_traits must specify chromosome key"})
				}
				pattern := filepath.Join(g.Root, "traits", "chromosome", name, "trait.yaml")
				if _, err := os.Stat(pattern); err != nil {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_trait_file_exists", Message: "chromosome trait file not found: " + pattern})
				}
			}
		}
	}
}
