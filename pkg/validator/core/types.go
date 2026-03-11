package core

import (
	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
)

// Severity indicates how a rule violation should be treated.
type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
)

// Issue is a single validation finding.
type Issue struct {
	Severity Severity
	Code     string
	Message  string
	Gene     string
	Codon    string
}

// Result aggregates issues produced by validation rules.
type Result struct {
	Issues []Issue
}

// Add appends an issue.
func (r *Result) Add(issue Issue) {
	r.Issues = append(r.Issues, issue)
}

// HasErrors returns true if any issue is an error.
func (r Result) HasErrors() bool {
	for _, is := range r.Issues {
		if is.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Summary reports counts per severity.
func (r Result) Summary() (errors, warns, infos int) {
	for _, is := range r.Issues {
		switch is.Severity {
		case SeverityError:
			errors++
		case SeverityWarn:
			warns++
		case SeverityInfo:
			infos++
		}
	}
	return
}

// Rule is a validation function.
type Rule func(g *loader.Genome, env map[string]nt.TypeNode, res *Result)

var registry []Rule

// Register adds a rule to the registry.
func Register(r Rule) {
	registry = append(registry, r)
}

// All returns registered rules.
func All() []Rule {
	return registry
}
