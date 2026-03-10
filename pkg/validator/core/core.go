package core

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
)

// Result holds validation findings.
type Result struct {
	Errors []string
}

func (r *Result) Add(path, msg string) {
	r.Errors = append(r.Errors, fmt.Sprintf("%s: %s", path, msg))
}

func (r *Result) Err() error {
	if len(r.Errors) == 0 {
		return nil
	}
	return fmt.Errorf("validation failed:\n- %s", strings.Join(r.Errors, "\n- "))
}

// Rule is a validation rule plug‑in.
type Rule func(*loader.ComposedGenome, *Result)

var registry []Rule

// Register adds a rule to the global registry. Called from rule files' init().
func Register(r Rule) {
	registry = append(registry, r)
}

// Registry returns all registered rules (copy).
func Registry() []Rule {
	out := make([]Rule, len(registry))
	copy(out, registry)
	return out
}
