package main

import (
	"fmt"

	"github.com/Spencer1O1/codon_language/pkg/loader"
)

// printLoaderIssues displays loader-stage issues and returns true if any errors were present.
func printLoaderIssues(issues []loader.Issue) bool {
	if len(issues) == 0 {
		return false
	}
	err := false
	for _, is := range issues {
		printIssue("loader", is.Severity, is.Code, is.Message, "", "")
		if is.Severity == "error" {
			err = true
		}
	}
	return err
}

// printIssue renders a single issue consistently for loader and validator.
func printIssue(prefix, severity, code, message, gene, codon string) {
	line := fmt.Sprintf("%s-%s [%s]: %s", prefix, severity, code, message)
	if gene != "" || codon != "" {
		line += fmt.Sprintf(" (gene=%s codon=%s)", gene, codon)
	}
	fmt.Println(line)
}
