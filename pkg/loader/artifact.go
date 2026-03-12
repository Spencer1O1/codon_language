package loader

import (
	"sort"
	"strings"

	tp "github.com/Spencer1O1/codon_language/pkg/nucleotype"
)

// ComposedArtifact is the serialized, post-validation genome shape.
type ComposedArtifact struct {
	SchemaVersion string               `yaml:"schema_version" json:"schema_version"`
	Manifest      map[string]any       `yaml:"manifest" json:"manifest"`
	CodonSchemas  map[string]any       `yaml:"codon_schemas" json:"codon_schemas"`
	Nucleotypes   map[string]string    `yaml:"nucleotypes" json:"nucleotypes"`
	Chromosomes   []ArtifactChromosome `yaml:"chromosomes" json:"chromosomes"`
	Expression    map[string]any       `yaml:"expression,omitempty" json:"expression,omitempty"`
	TraitsApplied []TraitApplied       `yaml:"traits_applied,omitempty" json:"traits_applied,omitempty"`
	Issues        []Issue              `yaml:"issues,omitempty" json:"issues,omitempty"`
}

// ArtifactChromosome groups genes by chromosome.
type ArtifactChromosome struct {
	Name  string         `yaml:"name" json:"name"`
	Genes []ArtifactGene `yaml:"genes" json:"genes"`
}

// ArtifactGene is a gene entry in the composed artifact.
type ArtifactGene struct {
	Name        string         `yaml:"name" json:"name"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Codons      map[string]any `yaml:"codons" json:"codons"`
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

	// group genes by chromosome for determinism
	chMap := map[string][]ArtifactGene{}
	for _, ge := range g.Genes {
		chMap[ge.Chromosome] = append(chMap[ge.Chromosome], ArtifactGene{Name: ge.Name, Description: ge.Description, Codons: ge.Codons})
	}
	var chromosomes []ArtifactChromosome
	for ch, genes := range chMap {
		sort.Slice(genes, func(i, j int) bool { return genes[i].Name < genes[j].Name })
		chromosomes = append(chromosomes, ArtifactChromosome{Name: ch, Genes: genes})
	}
	sort.Slice(chromosomes, func(i, j int) bool { return chromosomes[i].Name < chromosomes[j].Name })
	art.Chromosomes = chromosomes

	if g.Expression != nil {
		exp := map[string]any{}
		if g.Expression.Targets != nil {
			exp["targets"] = g.Expression.Targets
		}
		if g.Expression.Projections != nil {
			exp["projections"] = g.Expression.Projections
		}
		if g.Expression.Styles != nil {
			exp["styles"] = g.Expression.Styles
		}
		if g.Expression.Templates != nil {
			exp["templates"] = g.Expression.Templates
		}
		if len(exp) > 0 {
			art.Expression = exp
		}
	}

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
