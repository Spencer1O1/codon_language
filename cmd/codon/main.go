package main

import (
	"fmt"
	"os"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := args[1]
	root := ".codon"
	if len(args) > 2 {
		root = args[2]
	}

	switch cmd {
	case "load":
		if err := runLoad(root); err != nil {
			fmt.Fprintf(os.Stderr, "load error: %v\n", err)
			os.Exit(1)
		}
	case "validate":
		if err := runValidate(root); err != nil {
			fmt.Fprintf(os.Stderr, "validate error: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: codon <load|validate> [path]")
}

func runLoad(root string) error {
	g, err := loader.LoadGenome(root)
	if err != nil {
		return err
	}
	printLoaderIssues(g.Issues)
	fmt.Printf("codon schemas: %d, genes: %d\n", len(g.Schemas), len(g.Genes))
	return nil
}

func runValidate(root string) error {
	g, err := loader.LoadGenome(root)
	if err != nil {
		return err
	}
	if hasLoaderErrors := printLoaderIssues(g.Issues); hasLoaderErrors {
		return fmt.Errorf("loader reported errors")
	}
	env, err := loader.BuildTypeEnv(root)
	if err != nil {
		return err
	}
	res := validator.Validate(g, env)
	errs, warns, infos := res.Summary()
	fmt.Printf("errors: %d, warnings: %d, infos: %d\n", errs, warns, infos)
	for _, is := range res.Issues {
		fmt.Printf("%s [%s]: %s (gene=%s codon=%s)\n", is.Severity, is.Code, is.Message, is.Gene, is.Codon)
	}
	if res.HasErrors() {
		return fmt.Errorf("validation failed")
	}
	return nil
}

// printLoaderIssues displays loader-stage issues and returns true if any errors were present.
func printLoaderIssues(issues []loader.Issue) bool {
	if len(issues) == 0 {
		return false
	}
	errs := false
	for _, is := range issues {
		fmt.Printf("loader-%s [%s]: %s\n", is.Severity, is.Code, is.Message)
		if is.Severity == "error" {
			errs = true
		}
	}
	return errs
}
