package loader

import (
	"fmt"
	"io/fs"
	"os"
	path "path/filepath"
	"sort"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/core_assets"
	"github.com/Spencer1O1/codon-language/internal/engine_assets"
	tp "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	goyaml "gopkg.in/yaml.v3"
)

// Genome is a minimal composed genome for loader.
type Genome struct {
	Schemas      map[string]CodonSchema
	SchemaExport []string
	TypeEnv      map[string]tp.TypeNode
	TypeExport   []string
	Genes        []Gene
	Manifest     map[string]any
	Expression   *ExpressionAssets
	Root         string
	Issues       []Issue
}

type CodonSchema struct {
	Version     string
	Description string
	TypeExpr    string
	TypeAST     tp.TypeNode
	Source      string // filename
}

type Gene struct {
	Name        string
	Chromosome  string
	Description string
	Codons      map[string]any
	Path        string
}

type Issue struct {
	Severity string
	Code     string
	Message  string
}

// LoadGenome loads codon schemas and genes from a loader root.
func LoadGenome(root string) (*Genome, error) {
	codonSchemas, schemaExport, issues, err := loadCodonSchemas(root)
	if err != nil {
		return nil, err
	}
	manifest, err := loadManifest(root)
	if err != nil {
		return nil, err
	}
	expr, exprIssues := loadExpression(root)
	issues = append(issues, exprIssues...)
	genes, err := loadGenes(root)
	if err != nil {
		return nil, err
	}
	typeEnv, exportedTypes, err := BuildTypeEnv(root)
	if err != nil {
		return nil, err
	}
	return &Genome{
		Schemas:      codonSchemas,
		SchemaExport: schemaExport,
		TypeEnv:      typeEnv,
		TypeExport:   exportedTypes,
		Genes:        genes,
		Manifest:     manifest,
		Expression:   expr,
		Root:         root,
		Issues:       issues,
	}, nil
}

func loadCodonSchemas(root string) (map[string]CodonSchema, []string, []Issue, error) {
	codonSchemas := map[string]CodonSchema{}
	var exportSchemas []string
	var issues []Issue
	// embedded defaults
	if err := loadCodonSchemasFromFS(core_assets.CodonSchemas, "codon_schemas", codonSchemas, true, &exportSchemas); err != nil {
		return nil, exportSchemas, issues, err
	}
	// engine-only embedded
	if err := loadCodonSchemasFromFS(engine_assets.CodonSchemas, "codon_schemas", codonSchemas, false, nil); err != nil {
		return nil, exportSchemas, issues, err
	}
	// disk overrides/extensions
	files, err := path.Glob(path.Join(root, "codon_schemas", "*.yaml"))
	if err != nil {
		return nil, exportSchemas, issues, err
	}
	sort.Strings(files)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		names, err := parseSchemaDoc(data, f, codonSchemas)
		if err != nil {
			issues = append(issues, Issue{Severity: "error", Code: "schema_parseable", Message: err.Error()})
		} else {
			exportSchemas = append(exportSchemas, names...)
		}
	}
	return codonSchemas, exportSchemas, issues, nil
}

func loadManifest(root string) (map[string]any, error) {
	path := path.Join(root, "genome.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest map[string]any
	if err := goyaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	return manifest, nil
}

// ExpressionAssets holds parsed expression files (optional).
type ExpressionAssets struct {
	Targets     map[string]any
	Projections map[string]any
	Styles      map[string]any
	Templates   map[string]any
}

func loadExpression(root string) (*ExpressionAssets, []Issue) {
	exprRoot := path.Join(root, "expression")
	if _, err := os.Stat(exprRoot); err != nil {
		return nil, nil
	}
	var issues []Issue
	loadFile := func(name string) map[string]any {
		fp := path.Join(exprRoot, name)
		data, err := os.ReadFile(fp)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			issues = append(issues, Issue{Severity: "error", Code: "expression_read_failed", Message: fmt.Sprintf("read %s: %v", fp, err)})
			return nil
		}
		var m map[string]any
		if err := goyaml.Unmarshal(data, &m); err != nil {
			issues = append(issues, Issue{Severity: "error", Code: "expression_parse_failed", Message: fmt.Sprintf("parse %s: %v", fp, err)})
			return nil
		}
		return m
	}
	return &ExpressionAssets{
		Targets:     loadFile("targets.yaml"),
		Projections: loadFile("projections.yaml"),
		Styles:      loadFile("styles.yaml"),
		Templates:   loadFile("templates.yaml"),
	}, issues
}

func loadCodonSchemasFromFS(fsys fs.FS, dir string, dest map[string]CodonSchema, export bool, exportList *[]string) error {
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
		names, err := parseSchemaDoc(data, e.Name(), dest)
		if err != nil {
			return err
		}
		if export && exportList != nil {
			*exportList = append(*exportList, names...)
		}
	}
	return nil
}

func parseSchemaDoc(data []byte, filename string, dest map[string]CodonSchema) ([]string, error) {
	var doc struct {
		Codons map[string]struct {
			Version     string `yaml:"version"`
			Description string `yaml:"description"`
			Schema      string `yaml:"schema"`
			TypeLegacy  string `yaml:"type"`
		} `yaml:"codons"`
	}
	if err := goyaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse codon schema file %s: %w", filename, err)
	}
	var names []string
	for name, cf := range doc.Codons {
		src := cf.Schema
		if strings.TrimSpace(src) == "" {
			src = cf.TypeLegacy // backward compat
		}
		if strings.TrimSpace(src) == "" {
			continue
		}
		ast, err := tp.Parse(src)
		if err != nil {
			return nil, fmt.Errorf("parse codon schema %s type %s: %w", filename, name, err)
		}
		dest[name] = CodonSchema{Version: cf.Version, Description: cf.Description, TypeExpr: src, TypeAST: ast, Source: filename}
		names = append(names, name)
	}
	return names, nil
}

// ParseSchemaDocInto parses schemas into an existing map (used for trait-local schemas).
func ParseSchemaDocInto(data []byte, filename string, dest map[string]CodonSchema) error {
	_, err := parseSchemaDoc(data, filename, dest)
	return err
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
		// gene_required validation rule
		if !ok || name == "" {
			return nil, fmt.Errorf("gene file %s missing gene name", p)
		}
		desc, _ := raw["description"].(string)
		codons := map[string]any{}
		// codons_required validation rule
		if c, ok := raw["codons"].(map[string]any); ok {
			codons = c
		}
		chrom := chromosomeFromPath(root, p)
		genes = append(genes, Gene{Name: name, Chromosome: chrom, Description: desc, Codons: codons, Path: p})
	}
	return genes, nil
}

func chromosomeFromPath(root, full string) string {
	rel, err := path.Rel(path.Join(root, "chromosomes"), full)
	if err != nil {
		return ""
	}
	parts := strings.Split(rel, string(os.PathSeparator))
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// BuildTypeEnv builds a symbol table from nucleotide declarations.
// Returns the full env and the list of exported type names (public).
func BuildTypeEnv(root string) (map[string]tp.TypeNode, []string, error) {
	env := map[string]tp.TypeNode{}
	var exported []string
	// embedded defaults (public primitives)
	if err := loadTypesFromFS(core_assets.Nucleotypes, "nucleotides/types", env, true, &exported); err != nil {
		return nil, exported, err
	}
	// engine-only embedded types (not exported)
	if err := loadTypesFromFS(engine_assets.Nucleotypes, "nucleotides/types", env, false, nil); err != nil {
		return nil, exported, err
	}
	// hardcoded internal nucleotypes (not shipped to consumers)
	coreNames := `export chromosome_name = string
export gene_name = string
export genome_trait_name = string
export chromosome_trait_name = string
export gene_trait_name = string

export entity_name = string
export capability_name = string
export relation_name = string
`
	if _, err := parseTypesDoc(coreNames, "builtin:names", env); err != nil {
		return nil, exported, err
	}
	coreFields := `field_key = Regex<"^[a-z][a-z0-9_]*$">
export field<T> = Map<field_key, {
  type: T,
  default: Optional<T>,
  unique: Optional<boolean>,
  optional: Optional<boolean>,
}>
`
	if _, err := parseTypesDoc(coreFields, "builtin:fields", env); err != nil {
		return nil, exported, err
	}
	// disk overrides/extensions
	files, err := path.Glob(path.Join(root, "nucleotides", "types", "*.nucleotype"))
	if err != nil {
		return nil, exported, err
	}
	sort.Strings(files)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, exported, err
		}
		names, err := parseTypesDoc(string(data), f, env)
		if err != nil {
			return nil, exported, err
		}
		exported = append(exported, names...)
	}
	// inject primitives as terminals
	env["primitive"] = tp.NameType{Name: "primitive"}
	return env, exported, nil
}

func loadTypesFromFS(fsys fs.FS, dir string, env map[string]tp.TypeNode, export bool, exportList *[]string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		isPrimitive := e.Name() == "primitives.nucleotype"
		if export && !isPrimitive {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		names, err := parseTypesDoc(string(data), e.Name(), env)
		if err != nil {
			return err
		}
		if export && exportList != nil {
			*exportList = append(*exportList, names...)
		}
	}
	return nil
}

func parseTypesDoc(src string, filename string, env map[string]tp.TypeNode) ([]string, error) {
	decls, err := tp.ParseDecls(src)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}
	var names []string
	for _, d := range decls {
		env[d.Name] = d.Type
		names = append(names, d.Name)
	}
	return names, nil
}

// ParseTypesDoc exported for reuse (e.g., trait-local nucleotypes).
func ParseTypesDoc(src string, filename string, env map[string]tp.TypeNode) error {
	_, err := parseTypesDoc(src, filename, env)
	return err
}
