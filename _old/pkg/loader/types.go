package loader

import "encoding/json"

// ComposedGenome represents the composed_genome_contract described in
// .codon/chromosomes/genome/composition.yaml. It is the interoperability
// boundary handed to validators and expression engines.
type ComposedGenome struct {
	SchemaVersion string         `json:"schema_version" yaml:"schema_version"`
	Project       Project        `json:"project" yaml:"project"`
	Traits        []string       `json:"traits,omitempty" yaml:"traits,omitempty"`
	Genes         []ComposedGene `json:"genes" yaml:"genes"`
	CodonFamilies map[string]CodonFamily `json:"codon_families,omitempty" yaml:"codon_families,omitempty"`
}

// MarshalJSON keeps explicit aliasing for stable output. Ordering is applied
// earlier via deterministic sorting in Load.
func (c *ComposedGenome) MarshalJSON() ([]byte, error) {
	type alias ComposedGenome
	return json.Marshal((*alias)(c))
}

// ComposedGene is the normalized representation of a gene after loading.
type ComposedGene struct {
	Name         string                `json:"name" yaml:"name"`
	Chromosome   string                `json:"chromosome" yaml:"chromosome"`
	Purpose      string                `json:"purpose,omitempty" yaml:"purpose,omitempty"`
	Dependencies []string              `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Entities     []ComposedEntity      `json:"entities,omitempty" yaml:"entities,omitempty"`
	Capabilities []ComposedCapability  `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Relations    []RelationDefinition  `json:"relations,omitempty" yaml:"relations,omitempty"`
	References   []ReferenceDefinition `json:"references,omitempty" yaml:"references,omitempty"`
	Traits       []string              `json:"traits,omitempty" yaml:"traits,omitempty"`
	RawCodons    map[string]any        `json:"raw_codons,omitempty" yaml:"raw_codons,omitempty"`
}

// CodonFamily describes an extensible codon family declared in codon_families.yaml.
type CodonFamily struct {
	Name            string         `json:"name" yaml:"name"`
	Description     string         `json:"description,omitempty" yaml:"description,omitempty"`
	Version         string         `json:"version,omitempty" yaml:"version,omitempty"`
	SchemaRef       string         `json:"schema_ref,omitempty" yaml:"schema_ref,omitempty"`
	ProjectionHints map[string]any `json:"projection_hints,omitempty" yaml:"projection_hints,omitempty"`
}

// ComposedEntity is the normalized entity form.
type ComposedEntity struct {
	Name   string                     `json:"name" yaml:"name"`
	Fields map[string]FieldDefinition `json:"fields,omitempty" yaml:"fields,omitempty"`
}

// ComposedCapability is the normalized capability form. Both semantic
// (kebab-case string) and structured forms are folded into this shape.
type ComposedCapability struct {
	Name        string                     `json:"name" yaml:"name"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Inputs      map[string]FieldDefinition `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs     map[string]FieldDefinition `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Effects     []string                   `json:"effects,omitempty" yaml:"effects,omitempty"`
}

// Manifest mirrors the structure of genome.yaml used during loading.
type Manifest struct {
	SchemaVersion string         `json:"schema_version" yaml:"schema_version"`
	Project       Project        `json:"project" yaml:"project"`
	Traits        []string       `json:"traits,omitempty" yaml:"traits,omitempty"`
	Expression    map[string]any `json:"expression,omitempty" yaml:"expression,omitempty"`
}

// Project captures manifest.project.
type Project struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
}

// FieldDefinition mirrors language.field.field_definition.
type FieldDefinition struct {
	Type      string   `json:"type" yaml:"type"`
	Optional  bool     `json:"optional,omitempty" yaml:"optional,omitempty"`
	Unique    bool     `json:"unique,omitempty" yaml:"unique,omitempty"`
	Default   any      `json:"default,omitempty" yaml:"default,omitempty"`
	Values    []string `json:"values,omitempty" yaml:"values,omitempty"`
	Reference string   `json:"reference,omitempty" yaml:"reference,omitempty"`
}

// RelationDefinition mirrors language.relation.relation_definition.
type RelationDefinition struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
	Type string `json:"type" yaml:"type"`
	Name string `json:"name" yaml:"name"`
}

// ReferenceDefinition mirrors language.reference.reference_definition.
type ReferenceDefinition struct {
	From      string `json:"from" yaml:"from"`
	To        string `json:"to" yaml:"to"`
	Type      string `json:"type" yaml:"type"`
	Name      string `json:"name" yaml:"name"`
	Reference string `json:"reference" yaml:"reference"`
}
