package main

import (
	"fmt"

	"github.com/Spencer1O1/codon_language/pkg/loader"
	"github.com/Spencer1O1/codon_language/pkg/validator"
)

func runValidate(root string) error {
	g, err := loader.LoadGenome(root)
	if err != nil {
		return err
	}
	if hasLoaderErrors := printLoaderIssues(g.Issues); hasLoaderErrors {
		return fmt.Errorf("loader reported errors")
	}
	res := validator.Validate(g, g.TypeEnv)
	errs, warns, infos := res.Summary()
	fmt.Printf("errors: %d, warnings: %d, infos: %d\n", errs, warns, infos)
	for _, is := range res.Issues {
		printIssue("validation", string(is.Severity), is.Code, is.Message, is.Gene, is.Codon)
	}
	if res.HasErrors() {
		return fmt.Errorf("validation failed")
	}
	return nil
}
