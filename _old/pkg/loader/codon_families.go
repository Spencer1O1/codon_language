package loader

import (
	"os"

	"gopkg.in/yaml.v3"
)

// loadCodonFamilies loads optional standard/custom family registries from .codon/codons.
func loadCodonFamilies(root string) (map[string]CodonFamily, error) {
	paths := []string{
		root + "/codon_families/core.yaml",
		root + "/codon_families/standard.yaml",
		root + "/codon_families/custom.yaml",
	}
	families := map[string]CodonFamily{}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var raw struct {
			Families map[string]CodonFamily `yaml:"families"`
		}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		for name, cf := range raw.Families {
			if cf.Name == "" {
				cf.Name = name
			}
			families[name] = cf // custom overrides standard if same name
		}
	}
	if len(families) == 0 {
		return nil, nil
	}
	return families, nil
}
