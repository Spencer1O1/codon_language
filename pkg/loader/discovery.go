package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type geneFile struct {
	Chromosome string
	Path       string
}

func discoverGeneFiles(chromosomesDir string) ([]geneFile, error) {
	var files []geneFile
	entries, err := os.ReadDir(chromosomesDir)
	if err != nil {
		return nil, fmt.Errorf("list chromosomes: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		chromosome := entry.Name()
		if err := ValidateIdentifier("chromosome", chromosome); err != nil {
			return nil, fmt.Errorf("chromosome %q invalid: %w", chromosome, err)
		}
		chromoPath := filepath.Join(chromosomesDir, chromosome)
		children, err := os.ReadDir(chromoPath)
		if err != nil {
			return nil, fmt.Errorf("list genes in %s: %w", chromosome, err)
		}
		for _, child := range children {
			if child.IsDir() || strings.HasPrefix(child.Name(), ".") {
				continue
			}
			if filepath.Ext(child.Name()) != ".yaml" {
				continue
			}
			files = append(files, geneFile{
				Chromosome: chromosome,
				Path:       filepath.Join(chromoPath, child.Name()),
			})
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no gene files found under %s", chromosomesDir)
	}
	return files, nil
}
