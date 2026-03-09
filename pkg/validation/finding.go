package validation

type Severity string

const (
	SeverityError      Severity = "error"
	SeverityWarning    Severity = "warning"
	SeveritySuggestion Severity = "suggestion"
)

type Finding struct {
	Severity Severity
	Code     string
	Message  string
	Path     string
}

func (f Finding) IsError() bool {
	return f.Severity == SeverityError
}
