package language

import (
	"path/filepath"
	"strings"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
	"github.com/Spencer1O1/codon_language/pkg/validator/core"
)

func init() { core.RegisterWithGroup("language", codonSchemaRules) }

func codonSchemaRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for name, sc := range g.Schemas {
		if strings.TrimSpace(name) == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_name_required", Message: "codon schema name must be non-empty"})
		}
		if strings.TrimSpace(sc.Version) == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_version_required", Message: "codon schema version is required"})
		}
		if strings.TrimSpace(sc.TypeExpr) == "" {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_field_required", Message: "codon schema must include schema/type_expr"})
		}
		// filename extension only for disk-loaded schemas (Source has a path separator)
		if strings.Contains(sc.Source, string(filepath.Separator)) {
			if filepath.Ext(sc.Source) != ".yaml" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "schema_filename_extension", Message: "codon schema files must use .yaml extension"})
			}
		}
	}
}
