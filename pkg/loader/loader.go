package loader

import (
	fmt "fmt"
	os "os"
	path "path/filepath"
	sort "sort"
	strings "strings"

	tp "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	goyaml "gopkg.in/yaml.v3"
)

// Genome is a minimal composed genome for loader.
type Genome struct {
	Families map[string]Family
	Genes    []Gene
}

type Family struct {
	Version     string
	Description string
	TypeExpr    string
	TypeAST     tp.TypeNode
}

type Gene struct {
	Name    string
	Codons  map[string]any
	Imports []string
}

// LoadGenome loads families and genes from a loader root.
func LoadGenome(root string) (*Genome, error) {
	families, err := loadFamilies(root)
	if err != nil {
		return nil, err
	}
	genes, err := loadGenes(root)
	if err != nil {
		return nil, err
	}
	return &Genome{Families: families, Genes: genes}, nil
}

func loadFamilies(root string) (map[string]Family, error) {
	glob := path.Join(root, "codon_families", "*.codon")
	files, err := path.Glob(glob)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	families := map[string]Family{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var doc struct {
			Families map[string]struct {
				Version     string `yaml:"version"`
				Description string `yaml:"description"`
				Type        string `yaml:"type"`
			} `yaml:"families"`
		}
		if err := goyaml.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("parse family %s: %w", f, err)
		}
		for name, cf := range doc.Families {
			if strings.TrimSpace(cf.Type) == "" {
				continue
			}
			ast, err := tp.Parse(cf.Type)
			if err != nil {
				return nil, fmt.Errorf("parse family %s type %s: %w", f, name, err)
			}
			families[name] = Family{Version: cf.Version, Description: cf.Description, TypeExpr: cf.Type, TypeAST: ast}
		}
	}
	return families, nil
}

func loadGenes(root string) ([]Gene, error) {
	var genes []Gene
	glob := path.Join(root, "chromosomes", "**", "*.yaml")
	paths, err := path.Glob(glob)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var raw map[string]any
		if err := goyaml.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("parse gene %s: %w", p, err)
		}
		name, ok := raw["gene"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("gene file %s missing gene name", p)
		}
		imports := toStringList(raw["imports"])
		codons := map[string]any{}
		if c, ok := raw["codons"].(map[string]any); ok {
			codons = c
		}
		genes = append(genes, Gene{Name: name, Codons: codons, Imports: imports})
	}
	return genes, nil
}

func toStringList(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// BuildTypeEnv builds a symbol table from nucleotide declarations.
func BuildTypeEnv(root string) (map[string]tp.TypeNode, error) {
	env := map[string]tp.TypeNode{}
	files, err := path.Glob(path.Join(root, "nucleotides", "types", "*.nucleotype"))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		decls, err := tp.ParseDecls(string(data))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
		for _, d := range decls {
			ast := tp.Resolve(d.Type, env)
			env[d.Name] = ast
		}
	}
	// inject primitives as terminals
	env["primitive"] = tp.NameType{Name: "primitive"}
	env["any"] = tp.NameType{Name: "any"}
	return env, nil
}
