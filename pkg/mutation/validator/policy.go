package validator

// ShouldApplyStrict blocks on errors and, when blockOnWarnings is true, also blocks on warnings.
func ShouldApplyStrict(res *Result, blockOnWarnings bool) bool {
	for _, f := range res.Findings {
		if f.Severity == Error {
			return false
		}
		if blockOnWarnings && f.Severity == Warning {
			return false
		}
	}
	return true
}
