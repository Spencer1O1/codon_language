package validator

// ShouldApply reports whether a patch can proceed (no error severities).
func ShouldApply(res *Result) bool {
	for _, f := range res.Findings {
		if f.Severity == Error {
			return false
		}
	}
	return true
}
