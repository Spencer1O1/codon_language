package validation

import "fmt"

func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.IsError() {
			return true
		}
	}
	return false
}

func PrintFindings(title string, findings []Finding) {
	fmt.Printf("\n%s:\n", title)

	if len(findings) == 0 {
		fmt.Println("OK")
		return
	}

	for _, f := range findings {
		fmt.Printf("[%s] %s (%s)\n", f.Severity, f.Message, f.Path)
	}
}
