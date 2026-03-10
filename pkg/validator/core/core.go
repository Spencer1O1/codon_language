package core

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
)

// Result holds validation findings.
type Result struct {
	Findings []Finding
}

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Info    Severity = "info"
)

type Finding struct {
	Severity Severity
	Path     string
	Message  string
}

// Summary aggregates findings by severity.
type Summary struct {
	Errors   int
	Warnings int
	Infos    int
	Total    int
}

func (r *Result) Add(path, msg string) {
	r.AddWithSeverity(Error, path, msg)
}

func (r *Result) AddWithSeverity(sev Severity, path, msg string) {
	r.Findings = append(r.Findings, Finding{Severity: sev, Path: path, Message: msg})
}

func (r *Result) Err() error {
	if len(r.Findings) == 0 {
		return nil
	}
	var errs []string
	for _, f := range r.Findings {
		if f.Severity == Error {
			errs = append(errs, fmt.Sprintf("%s: %s", f.Path, f.Message))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("validation failed:\n- %s", strings.Join(errs, "\n- "))
}

// Summary returns counts of findings per severity.
func (r *Result) Summary() Summary {
	var s Summary
	for _, f := range r.Findings {
		s.Total++
		switch f.Severity {
		case Error:
			s.Errors++
		case Warning:
			s.Warnings++
		case Info:
			s.Infos++
		}
	}
	return s
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
