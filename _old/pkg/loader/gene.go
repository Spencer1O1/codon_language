package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func loadGenes(files []geneFile) ([]ComposedGene, error) {
	var genes []ComposedGene
	for _, gf := range files {
		gene, err := loadGeneFile(gf)
		if err != nil {
			return nil, err
		}
		genes = append(genes, gene)
	}
	return genes, nil
}

func loadGeneFile(gf geneFile) (ComposedGene, error) {
	content, err := os.ReadFile(gf.Path)
	if err != nil {
		return ComposedGene{}, fmt.Errorf("read gene file %s: %w", gf.Path, err)
	}
	var raw map[string]any
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return ComposedGene{}, fmt.Errorf("parse gene file %s: %w", gf.Path, err)
	}
	if err := expectKeys(gf.Path, raw, []string{"gene", "purpose", "types", "dependencies", "codons"}); err != nil {
		return ComposedGene{}, err
	}

	name, err := parseGeneName(gf.Path, raw)
	if err != nil {
		return ComposedGene{}, err
	}
	purpose, err := parsePurpose(gf.Path, raw)
	if err != nil {
		return ComposedGene{}, err
	}
	deps, err := parseDependencies(gf.Path, raw)
	if err != nil {
		return ComposedGene{}, err
	}
	entities, capabilities, relations, references, traits, rawCodons, err := parseCodons(gf.Path, raw)
	if err != nil {
		return ComposedGene{}, err
	}

	return ComposedGene{
		Name:         name,
		Chromosome:   gf.Chromosome,
		Purpose:      purpose,
		Dependencies: deps,
		Entities:     entities,
		Capabilities: capabilities,
		Relations:    relations,
		References:   references,
		Traits:       traits,
		RawCodons:    rawCodons,
	}, nil
}

func parseGeneName(path string, raw map[string]any) (string, error) {
	geneBlock, ok := raw["gene"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("%s: gene must be an object", path)
	}
	if err := expectKeys("gene", geneBlock, []string{"name"}); err != nil {
		return "", err
	}
	name, ok := geneBlock["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("%s: gene.name must be a non-empty string", path)
	}
	if err := ValidateIdentifier("gene", name); err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	fileStem := strings.TrimSuffix(filepath.Base(path), ".yaml")
	if name != fileStem {
		return "", fmt.Errorf("%s: gene.name %q must match filename %q", path, name, fileStem)
	}
	return name, nil
}

func parsePurpose(path string, raw map[string]any) (string, error) {
	if p, ok := raw["purpose"]; ok {
		if s, ok := p.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("%s: purpose must be a string when present", path)
	}
	return "", nil
}

func parseDependencies(path string, raw map[string]any) ([]string, error) {
	if depsRaw, ok := raw["dependencies"]; ok {
		items, err := toStringList("dependencies", depsRaw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		for _, d := range items {
			if err := ValidateGeneReference(d); err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
		}
		return items, nil
	}
	return nil, nil
}

func parseCodons(path string, raw map[string]any) ([]ComposedEntity, []ComposedCapability, []RelationDefinition, []ReferenceDefinition, []string, map[string]any, error) {
	if rawCodons, ok := raw["codons"]; ok {
		codons, ok := rawCodons.(map[string]any)
		if !ok {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("%s: codons must be an object", path)
		}
		rawCopy := make(map[string]any, len(codons))
		for k, v := range codons {
			rawCopy[k] = v
		}
		entities, err := parseEntities(codons["entities"])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("%s: %w", path, err)
		}
		capabilities, err := parseCapabilities(codons["capabilities"])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("%s: %w", path, err)
		}
		relations, err := parseRelations(codons["relations"])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("%s: %w", path, err)
		}
		references, err := parseReferences(codons["references"])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("%s: %w", path, err)
		}
		traits := parseTraitsList(codons["traits"])
		return entities, capabilities, relations, references, traits, rawCopy, nil
	}
	return nil, nil, nil, nil, nil, nil, nil
}
