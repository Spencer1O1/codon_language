package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Load composes a genome rooted at path and returns the normalized contract.
// It enforces key structural rules from composition.yaml, manifest.yaml,
// addressing.yaml, and gene.yaml. Trait expansion and deep semantic validation
// are left for later passes.
func Load(root string) (*ComposedGenome, error) {
	manifest, err := loadManifest(filepath.Join(root, "genome.yaml"))
	if err != nil {
		return nil, err
	}

	codonFamilies, err := loadCodonFamilies(root)
	if err != nil {
		return nil, fmt.Errorf("load codon_families.yaml: %w", err)
	}

	chromosomesDir := filepath.Join(root, "chromosomes")
	if err := ensureDir(chromosomesDir); err != nil {
		return nil, err
	}

	geneFiles, err := discoverGeneFiles(chromosomesDir)
	if err != nil {
		return nil, err
	}

	genes, err := loadGenes(geneFiles)
	if err != nil {
		return nil, err
	}

	orderGenes(genes)

	return &ComposedGenome{
		SchemaVersion: manifest.SchemaVersion,
		Project:       manifest.Project,
		Traits:        manifest.Traits,
		Genes:         genes,
		CodonFamilies: codonFamilies,
	}, nil
}

func ensureDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("missing required directory %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s must be a directory", path)
	}
	return nil
}

func orderGenes(genes []ComposedGene) {
	sort.Slice(genes, func(i, j int) bool {
		if genes[i].Chromosome == genes[j].Chromosome {
			return genes[i].Name < genes[j].Name
		}
		return genes[i].Chromosome < genes[j].Chromosome
	})
	for idx := range genes {
		sort.Slice(genes[idx].Entities, func(a, b int) bool {
			return genes[idx].Entities[a].Name < genes[idx].Entities[b].Name
		})
		sort.Slice(genes[idx].Capabilities, func(a, b int) bool {
			return genes[idx].Capabilities[a].Name < genes[idx].Capabilities[b].Name
		})
	}
}
