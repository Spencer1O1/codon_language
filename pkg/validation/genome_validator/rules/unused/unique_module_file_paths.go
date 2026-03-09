package unused

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

// NOT USED BECAUSE LOADER ALREADY ENFORCES THIS

type UniqueModuleFilePathsRule struct{}

func (UniqueModuleFilePathsRule) Name() string {
	return "unique-module-file-paths"
}

func (UniqueModuleFilePathsRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	pathToModules := make(map[string][]string)

	for moduleName, path := range g.Manifest.Modules {
		pathToModules[path] = append(pathToModules[path], moduleName)
	}

	for path, modules := range pathToModules {
		if len(modules) <= 1 {
			continue
		}

		sort.Strings(modules)

		findings = append(findings, validation.Finding{
			Severity: validation.SeverityError,
			Code:     "duplicate_module_file_path",
			Path:     "/modules",
			Message: fmt.Sprintf(
				"multiple modules reference the same module file path %q: %s",
				path,
				strings.Join(modules, ", "),
			),
		})
	}

	return findings
}
