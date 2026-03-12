package core

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
