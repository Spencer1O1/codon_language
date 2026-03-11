package validator

import (
	"fmt"
	"os"
	path "path/filepath"
	"reflect"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
	goyaml "gopkg.in/yaml.v3"
)

// applyTraitInjection mutates the genome by injecting genome and gene traits per spec.
func applyTraitInjection(g *loader.Genome, res *core.Result) {
	applyGenomeTraits(g, res)
	applyGeneTraits(g, res)
}

type genomeTraitFile struct {
	Genes map[string]map[string]any `yaml:"genes"`
}

type geneTraitFile struct {
	Codons map[string]any `yaml:"codons"`
}

func applyGenomeTraits(g *loader.Genome, res *core.Result) {
	traitsRaw, ok := g.Manifest["traits"]
	if !ok {
		return
	}
	traits, ok := traitsRaw.(map[string]any)
	if !ok {
		return
	}

	for chrom, val := range traits {
		name, ok := val.(string)
		if !ok {
			continue
		}
		tf, err := loadGenomeTraitFile(path.Join(g.Root, "traits", "genome", name+".yaml"))
		if err != nil {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_file_exists", Message: err.Error()})
			continue
		}
		injectedSeen := map[string]map[string]any{}
		for geneName, codons := range tf.Genes {
			genePtr := findGenePtr(g.Genes, chrom, geneName)
			if genePtr == nil {
				// create new gene
				g.Genes = append(g.Genes, loader.Gene{Chromosome: chrom, Name: geneName, Codons: deepCopyMap(codons)})
				genePtr = &g.Genes[len(g.Genes)-1]
			} else {
				mergeCodons(genePtr, codons, injectedSeen, res)
			}
		}
	}
}

func applyGeneTraits(g *loader.Genome, res *core.Result) {
	for gi := range g.Genes {
		gene := &g.Genes[gi]
		traitsRaw, ok := gene.Codons["traits"]
		if !ok {
			continue
		}
		list, ok := traitsRaw.([]any)
		if !ok {
			continue
		}
		injectedSeen := map[string]map[string]any{}
		for _, tr := range list {
			name, ok := tr.(string)
			if !ok {
				continue
			}
			tf, err := loadGeneTraitFile(path.Join(g.Root, "traits", "gene", name+".yaml"))
			if err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "gene_trait_file_exists", Message: err.Error(), Gene: gene.Name, Codon: "traits"})
				continue
			}
			mergeCodons(gene, tf.Codons, injectedSeen, res)
		}
	}
}

func mergeCodons(gene *loader.Gene, injected map[string]any, injectedSeen map[string]map[string]any, res *core.Result) {
	if injectedSeen == nil {
		injectedSeen = map[string]map[string]any{}
	}
	for codonName, val := range injected {
		if authored, ok := gene.Codons[codonName]; ok {
			if !reflect.DeepEqual(authored, val) {
				res.Add(core.Issue{Severity: core.SeverityWarn, Code: "trait_conflict_authored_wins", Message: fmt.Sprintf("codon %s on gene %s overridden by authored value", codonName, gene.Name), Gene: gene.Name, Codon: codonName})
			}
			continue
		}
		if prev, ok := injectedSeen[codonName]; ok {
			if !reflect.DeepEqual(prev, val) {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "trait_conflict_injected_must_match", Message: fmt.Sprintf("conflicting injected codon %s on gene %s", codonName, gene.Name), Gene: gene.Name, Codon: codonName})
				continue
			}
		}
		copied := deepCopyMapAny(val)
		injectedSeen[codonName] = asMapAny(copied)
		gene.Codons[codonName] = copied
	}
}

func deepCopyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = deepCopyMapAny(v)
	}
	return dst
}

func deepCopyMapAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return deepCopyMap(t)
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = deepCopyMapAny(e)
		}
		return out
	default:
		return t
	}
}

func asMapAny(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}

func loadGenomeTraitFile(path string) (*genomeTraitFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tf genomeTraitFile
	if err := goyaml.Unmarshal(data, &tf); err != nil {
		return nil, err
	}
	return &tf, nil
}

func loadGeneTraitFile(path string) (*geneTraitFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tf geneTraitFile
	if err := goyaml.Unmarshal(data, &tf); err != nil {
		return nil, err
	}
	// if no codons key, treat whole doc as codons map
	if len(tf.Codons) == 0 {
		var anyMap map[string]any
		if err := goyaml.Unmarshal(data, &anyMap); err == nil {
			tf.Codons = anyMap
		}
	}
	return &tf, nil
}

func findGenePtr(genes []loader.Gene, chrom, name string) *loader.Gene {
	for i := range genes {
		if genes[i].Chromosome == chrom && genes[i].Name == name {
			return &genes[i]
		}
	}
	return nil
}
