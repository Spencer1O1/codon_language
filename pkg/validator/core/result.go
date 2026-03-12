package core

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
