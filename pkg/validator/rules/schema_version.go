package rules

import (
	"fmt"
	"regexp"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

var semverRe = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

func init() {
	core.Register(checkSchemaVersion)
}

func checkSchemaVersion(genome *loader.ComposedGenome, res *core.Result) {
	if genome.SchemaVersion == "" {
		res.Add("genome.schema_version", "missing schema_version")
		return
	}
	if !semverRe.MatchString(genome.SchemaVersion) {
		res.Add("genome.schema_version", fmt.Sprintf("must follow semantic versioning, got %q", genome.SchemaVersion))
	}
}
