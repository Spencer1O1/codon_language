package validator

import (
	"fmt"
	"strings"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Info    Severity = "info"
)

// Finding represents a single validation notice.
type Finding struct {
	Severity Severity
	Code     string
	Path     string
	Message  string
}

// Result accumulates findings and computes aggregates.
type Result struct {
	Findings []Finding
}

func (r *Result) add(code string, sev Severity, path, msg string) {
	r.Findings = append(r.Findings, Finding{Severity: sev, Code: code, Path: path, Message: msg})
}

// Err returns an error when any error-level finding exists.
func (r *Result) Err() error {
	var errs []string
	for _, f := range r.Findings {
		if f.Severity == Error {
			errs = append(errs, fmt.Sprintf("%s: %s", f.Path, f.Message))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("patch validation failed:\n- %s", strings.Join(errs, "\n- "))
}

// HighestSeverity reports the highest encountered severity.
func (r *Result) HighestSeverity() Severity {
	highest := Info
	for _, f := range r.Findings {
		switch f.Severity {
		case Error:
			return Error
		case Warning:
			highest = Warning
		}
	}
	return highest
}

// Summary counts findings per severity.
type Summary struct {
	Errors   int
	Warnings int
	Infos    int
	Total    int
}

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
