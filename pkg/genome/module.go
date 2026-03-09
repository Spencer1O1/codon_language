package genome

type FieldValue any
type CapabilityValue any
type ServiceValue any

type Relation struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

type Reference struct {
	From      string `yaml:"from"`
	To        string `yaml:"to"`
	Type      string `yaml:"type"`
	Name      string `yaml:"name"`
	Reference string `yaml:"reference"`
}

type Entity struct {
	Fields map[string]FieldValue `yaml:"fields"`
}

type Module struct {
	Description  string                  `yaml:"description"`
	Dependencies []string                `yaml:"dependencies"`
	Traits       []string                `yaml:"traits"`
	Entities     map[string]Entity       `yaml:"entities"`
	Capabilities []CapabilityValue       `yaml:"capabilities"`
	Services     map[string]ServiceValue `yaml:"services"`
	Relations    []Relation              `yaml:"relations"`
	References   []Reference             `yaml:"references"`
}

type ModuleFile struct {
	Module map[string]Module `yaml:"module"`
}

type Genome struct {
	Manifest Manifest
	Modules  map[string]Module
}
