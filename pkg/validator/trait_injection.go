package validator

import (
	"fmt"
	"os"
	path "path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	tp "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
	goyaml "gopkg.in/yaml.v3"
)

// applyTraitInjection mutates the genome by injecting genome and gene traits per spec.
func applyTraitInjection(g *loader.Genome, res *core.Result) {
	applyGenomeTraits(g, res)
	applyChromosomeTraits(g, res)
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
		name = strings.TrimSuffix(path.Base(name), ".yaml")
			tf, schemas, err := loadTraitFile(path.Join(g.Root, "traits", "genome", name))
			if err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "genome_trait_file_exists", Message: err.Error()})
				continue
			}
			if err := mergeNucleotypesWithEnv(g, path.Join(g.Root, "traits", "genome", name), res); err != nil {
				// issues recorded; continue
			}
		if err := mergeSchemas(g, schemas, res); err != nil {
			// mergeSchemas reports issues; continue to attempt codon injection
		}
		injectedSeen := map[string]map[string]any{}
		for geneName, codons := range tf.Genes {
			genePtr := findGenePtr(g.Genes, chrom, geneName)
			if genePtr == nil {
				g.Genes = append(g.Genes, loader.Gene{Chromosome: chrom, Name: geneName, Codons: deepCopyMap(codons)})
				genePtr = &g.Genes[len(g.Genes)-1]
			} else {
				mergeCodons(genePtr, codons, injectedSeen, res)
			}
		}
	}
}

// applyChromosomeTraits applies traits scoped to a chromosome (manifest: chromosome_traits map).
func applyChromosomeTraits(g *loader.Genome, res *core.Result) {
	traitsRaw, ok := g.Manifest["chromosome_traits"]
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
		name = strings.TrimSuffix(path.Base(name), ".yaml")
		tp, schemas, err := loadTraitFile(path.Join(g.Root, "traits", "chromosome", name))
		if err != nil {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "chromosome_trait_file_exists", Message: err.Error()})
			continue
		}
		if err := mergeNucleotypesWithEnv(g, path.Join(g.Root, "traits", "chromosome", name), res); err != nil {
			// continue
		}
		if err := mergeSchemas(g, schemas, res); err != nil {
			// continue
		}
		injectedSeen := map[string]map[string]any{}
		for geneName, codons := range tp.Genes {
			genePtr := findGenePtr(g.Genes, chrom, geneName)
			if genePtr == nil {
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
			name = strings.TrimSuffix(path.Base(name), ".yaml")
			tf, schemas, err := loadTraitFile(path.Join(g.Root, "traits", "gene", name))
			if err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "gene_trait_file_exists", Message: err.Error(), Gene: gene.Name, Codon: "traits"})
				continue
			}
				if err := mergeNucleotypesWithEnv(g, path.Join(g.Root, "traits", "gene", name), res); err != nil {
					// issues recorded; continue
				}
			if err := mergeSchemas(g, schemas, res); err != nil {
				// issues recorded; continue
			}
			codons := tf.Genes[gene.Name]
			if codons == nil {
				codons = tf.Genes[""]
			}
			mergeCodons(gene, codons, injectedSeen, res)
		}
	}
}

func mergeCodons(gene *loader.Gene, injected map[string]any, injectedSeen map[string]map[string]any, res *core.Result) {
	if injectedSeen == nil {
		injectedSeen = map[string]map[string]any{}
	}
	for codonName, val := range injected {
		if authored, ok := gene.Codons[codonName]; ok {
			merged, issue := mergeValue(authored, val, gene.Name, codonName)
			if issue != nil {
				res.Add(*issue)
				continue
			}
			gene.Codons[codonName] = merged
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

// mergeValue applies trait merge policy:
// - maps: additive deep merge; conflicts on scalar-vs-map emit error and keep authored.
// - lists: append injected elements; no de-dupe.
// - scalars: authored wins; warn.
func mergeValue(authored, injected any, geneName, codon string) (any, *core.Issue) {
	switch a := authored.(type) {
	case map[string]any:
		b, ok := injected.(map[string]any)
		if !ok {
			return authored, &core.Issue{Severity: core.SeverityError, Code: "trait_shape_conflict", Message: fmt.Sprintf("codon %s on gene %s: authored map vs injected non-map", codon, geneName), Gene: geneName, Codon: codon}
		}
		return mergeMaps(a, b, geneName, codon)
	case []any:
		if b, ok := injected.([]any); ok {
			// append; we could dedupe but keep simple
			return append(deepCopySlice(a), deepCopySlice(b)...), nil
		}
		return authored, &core.Issue{Severity: core.SeverityError, Code: "trait_shape_conflict", Message: fmt.Sprintf("codon %s on gene %s: authored list vs injected non-list", codon, geneName), Gene: geneName, Codon: codon}
	default:
		// scalar authored wins
		if reflect.DeepEqual(authored, injected) {
			return authored, nil
		}
		return authored, &core.Issue{Severity: core.SeverityWarn, Code: "trait_conflict_authored_wins", Message: fmt.Sprintf("codon %s on gene %s: injected value ignored; authored wins", codon, geneName), Gene: geneName, Codon: codon}
	}
}

func mergeMaps(a, b map[string]any, geneName, codon string) (map[string]any, *core.Issue) {
	out := deepCopyMap(a)
	// stable order for determinism in tests
	keys := make([]string, 0, len(b))
	for k := range b {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if av, exists := out[k]; exists {
			mv, issue := mergeValue(av, b[k], geneName, fmt.Sprintf("%s.%s", codon, k))
			if issue != nil {
				return out, issue
			}
			out[k] = mv
		} else {
			out[k] = deepCopyMapAny(b[k])
		}
	}
	return out, nil
}

func deepCopySlice(in []any) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = deepCopyMapAny(v)
	}
	return out
}

func findGenePtr(genes []loader.Gene, chrom, name string) *loader.Gene {
	for i := range genes {
		if genes[i].Chromosome == chrom && genes[i].Name == name {
			return &genes[i]
		}
	}
	return nil
}

// loadTraitFile loads a genome or gene trait and any co-located custom schemas.
// traitPath may be a directory (containing <name>.yaml and optional custom_schemas.yaml)
// or a legacy flat file path stem.
func loadTraitFile(traitPath string) (*genomeTraitFile, map[string]loader.CodonSchema, error) {
	traitPath = strings.TrimSuffix(traitPath, ".yaml")
	// if traitPath is a directory, expect trait.yaml inside
	if info, err := os.Stat(traitPath); err == nil && info.IsDir() {
		file := path.Join(traitPath, "trait.yaml")
		if info2, err2 := os.Stat(file); err2 == nil && !info2.IsDir() {
			return loadTraitFileFrom(file, traitPath)
		}
		return nil, nil, fmt.Errorf("trait file not found: %s", file)
	}
	return nil, nil, fmt.Errorf("trait file not found: %s", traitPath)
}

func loadTraitFileFrom(file, dir string) (*genomeTraitFile, map[string]loader.CodonSchema, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	var tf genomeTraitFile
	if err := goyaml.Unmarshal(data, &tf); err != nil {
		return nil, nil, err
	}
	if len(tf.Genes) == 0 {
		// maybe it's a gene trait (flat codons)
		var gt geneTraitFile
		if err := goyaml.Unmarshal(data, &gt); err == nil && len(gt.Codons) > 0 {
			tf = genomeTraitFile{Genes: map[string]map[string]any{
				"": gt.Codons,
			}}
		}
	}
	schemas := map[string]loader.CodonSchema{}
	// load optional custom schemas
	customPath := path.Join(dir, "custom_schemas.yaml")
	if b, err := os.ReadFile(customPath); err == nil {
		if err := loader.ParseSchemaDocInto(b, customPath, schemas); err != nil {
			return &tf, schemas, err
		}
	}
	// also support codon_schemas/*.yaml inside trait dir
	if entries, err := os.ReadDir(path.Join(dir, "codon_schemas")); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			b, err := os.ReadFile(path.Join(dir, "codon_schemas", e.Name()))
			if err != nil {
				return &tf, schemas, err
			}
			if err := loader.ParseSchemaDocInto(b, e.Name(), schemas); err != nil {
				return &tf, schemas, err
			}
		}
	}
	return &tf, schemas, nil
}

// mergeSchemas adds new schemas to the genome, detecting conflicts.
func mergeSchemas(g *loader.Genome, add map[string]loader.CodonSchema, res *core.Result) error {
	for name, sc := range add {
		if existing, ok := g.Schemas[name]; ok {
			if existing.TypeExpr != sc.TypeExpr || existing.Version != sc.Version {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_conflict", Message: fmt.Sprintf("schema %s conflicts with existing definition", name)})
				continue
			}
			continue // identical is fine
		}
		g.Schemas[name] = sc
	}
	return nil
}

// mergeNucleotypes loads trait-local nucleotypes into the type env with conflict checks.
func mergeNucleotypesWithEnv(g *loader.Genome, traitDir string, res *core.Result) error {
	dir := traitDir
	if info, err := os.Stat(traitDir); err == nil && info.IsDir() {
		dir = traitDir
	} else {
		dir = path.Dir(traitDir)
	}
	ntDir := path.Join(dir, "nucleotides", "types")
	entries, err := os.ReadDir(ntDir)
	if err != nil {
		return nil // no local nucleotypes; ok
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(path.Join(ntDir, e.Name()))
		if err != nil {
			return err
		}
		localEnv := map[string]tp.TypeNode{}
		if err := loader.ParseTypesDoc(string(data), e.Name(), localEnv); err != nil {
			return err
		}
		for name, t := range localEnv {
			if existing, ok := g.TypeEnv[name]; ok {
				if !tp.Equal(existing, t) {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "nucleotype_conflict", Message: fmt.Sprintf("nucleotype %s conflicts with existing definition", name)})
					continue
				}
			} else {
				g.TypeEnv[name] = t
			}
		}
	}
	return nil
}
