package loader

import (
	"sort"
	"strings"

	tp "github.com/Spencer1O1/codon-language/pkg/nucleotype"
)

// ComposedArtifact is the serialized, post-validation genome shape.
type ComposedArtifact struct {
	SchemaVersion string            `yaml:"schema_version" json:"schema_version"`
	Manifest      map[string]any    `yaml:"manifest" json:"manifest"`
	CodonSchemas  map[string]any    `yaml:"codon_schemas" json:"codon_schemas"`
	Nucleotypes   map[string]string `yaml:"nucleotypes" json:"nucleotypes"`
	Genes         []ArtifactGene    `yaml:"genes" json:"genes"`
	TraitsApplied []TraitApplied    `yaml:"traits_applied,omitempty" json:"traits_applied,omitempty"`
	Issues        []Issue           `yaml:"issues,omitempty" json:"issues,omitempty"`
}

// ArtifactGene is a gene entry in the composed artifact.
type ArtifactGene struct {
	Chromosome string         `yaml:"chromosome" json:"chromosome"`
	Gene       string         `yaml:"gene" json:"gene"`
	Codons     map[string]any `yaml:"codons" json:"codons"`
}

// TraitApplied records a trait that was applied during composition.
type TraitApplied struct {
	Scope  string `yaml:"scope" json:"scope"`
	Name   string `yaml:"name" json:"name"`
	Target string `yaml:"target,omitempty" json:"target,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty"`
}

// BuildArtifact constructs a composed artifact from a loaded genome.
func BuildArtifact(g *Genome) *ComposedArtifact {
	art := &ComposedArtifact{
		Manifest:     g.Manifest,
		CodonSchemas: map[string]any{},
		Nucleotypes:  map[string]string{},
		Issues:       g.Issues,
	}

	// schema_version from manifest if present
	if v, ok := g.Manifest["schema_version"].(string); ok {
		art.SchemaVersion = v
	}

	// codon schemas: expose version/description/schema
	for _, name := range g.SchemaExport {
		cs := g.Schemas[name]
		art.CodonSchemas[name] = map[string]any{
			"version":     cs.Version,
			"description": cs.Description,
			"schema":      cs.TypeExpr,
			"source":      cs.Source,
		}
	}

	// nucleotypes: render TypeEnv to strings
	names := make([]string, 0, len(g.TypeExport))
	seen := map[string]bool{}
	for _, n := range g.TypeExport {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	sort.Strings(names)
	for _, n := range names {
		art.Nucleotypes[n] = formatType(g.TypeEnv[n])
	}

	// genes: sorted for determinism
	genes := make([]ArtifactGene, 0, len(g.Genes))
	for _, ge := range g.Genes {
		genes = append(genes, ArtifactGene{Chromosome: ge.Chromosome, Gene: ge.Name, Codons: ge.Codons})
	}
	sort.Slice(genes, func(i, j int) bool {
		if genes[i].Chromosome == genes[j].Chromosome {
			return genes[i].Gene < genes[j].Gene
		}
		return genes[i].Chromosome < genes[j].Chromosome
	})
	art.Genes = genes

	return art
}

// formatType renders a TypeNode back to a readable TypeExpr string.
func formatType(t tp.TypeNode) string {
	switch v := t.(type) {
	case tp.NameType:
		return v.Name
	case tp.LiteralType:
		return "\"" + v.Value + "\""
	case tp.GenericType:
		args := make([]string, 0, len(v.Args))
		for _, a := range v.Args {
			args = append(args, formatType(a))
		}
		return v.Name + "<" + strings.Join(args, ", ") + ">"
	case tp.OptionalType:
		return formatType(v.Base) + "?"
	case tp.ListType:
		return formatType(v.Base) + "[]"
	case tp.UnionType:
		parts := make([]string, 0, len(v.Options))
		for _, o := range v.Options {
			parts = append(parts, formatType(o))
		}
		return strings.Join(parts, " | ")
	case tp.ObjectType:
		// sort fields for determinism
		fields := make([]tp.Field, len(v.Fields))
		copy(fields, v.Fields)
		sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
		parts := make([]string, 0, len(fields))
		for _, f := range fields {
			parts = append(parts, f.Name+": "+formatType(f.Type))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return ""
	}
}
