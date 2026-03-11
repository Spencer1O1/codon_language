package loader

import (
	"fmt"
	"io/fs"
	"os"
	path "path/filepath"
	"sort"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/core_assets"
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
	families := map[string]Family{}
	// embedded defaults
	if err := loadFamiliesFromFS(core_assets.Families, "codon_schemas", families); err != nil {
		return nil, err
	}
	// disk overrides/extensions
	files, err := path.Glob(path.Join(root, "codon_schemas", "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		if err := parseFamilyDoc(data, f, families); err != nil {
			return nil, err
		}
	}
	return families, nil
}

func loadFamiliesFromFS(fsys fs.FS, dir string, dest map[string]Family) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		// if directory missing in embedded fs, treat as empty
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if err := parseFamilyDoc(data, e.Name(), dest); err != nil {
			return err
		}
	}
	return nil
}

func parseFamilyDoc(data []byte, filename string, dest map[string]Family) error {
	var doc struct {
		Families map[string]struct {
			Version     string `yaml:"version"`
			Description string `yaml:"description"`
			Schema      string `yaml:"schema"`
			TypeLegacy  string `yaml:"type"`
		} `yaml:"families"`
	}
	if err := goyaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse family %s: %w", filename, err)
	}
	for name, cf := range doc.Families {
		src := cf.Schema
		if strings.TrimSpace(src) == "" {
			src = cf.TypeLegacy // backward compat
		}
		if strings.TrimSpace(src) == "" {
			continue
		}
		ast, err := tp.Parse(src)
		if err != nil {
			return fmt.Errorf("parse family %s type %s: %w", filename, name, err)
		}
		dest[name] = Family{Version: cf.Version, Description: cf.Description, TypeExpr: src, TypeAST: ast}
	}
	return nil
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
	// embedded defaults
	if err := loadTypesFromFS(core_assets.Nucleotypes, "nucleotides/types", env); err != nil {
		return nil, err
	}
	// disk overrides/extensions
	files, err := path.Glob(path.Join(root, "nucleotides", "types", "*.nucleotype"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if err := parseTypesDoc(string(data), f, env); err != nil {
			return nil, err
		}
	}
	// inject primitives as terminals
	env["primitive"] = tp.NameType{Name: "primitive"}
	env["any"] = tp.NameType{Name: "any"}
	return env, nil
}

func loadTypesFromFS(fsys fs.FS, dir string, env map[string]tp.TypeNode) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if err := parseTypesDoc(string(data), e.Name(), env); err != nil {
			return err
		}
	}
	return nil
}

func parseTypesDoc(src string, filename string, env map[string]tp.TypeNode) error {
	decls, err := tp.ParseDecls(src)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filename, err)
	}
	for _, d := range decls {
		ast := tp.Resolve(d.Type, env)
		env[d.Name] = ast
	}
	return nil
}
