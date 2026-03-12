package main

import (
	"fmt"
	"os"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator"
	"gopkg.in/yaml.v3"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := args[1]
	root := ""
	if len(args) > 2 {
		root = args[2]
	}
	if root == "" {
		fmt.Fprintf(os.Stderr, "command syntax error: %v\n", "Genome root not specified.")
		os.Exit(1)
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
	case "emit":
		if err := runEmit(root); err != nil {
			fmt.Fprintf(os.Stderr, "emit error: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: codon <load|validate|emit> [path]")
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

// runEmit emits the composed genome artifact after successful validation.
func runEmit(root string) error {
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
	artifact := loader.BuildArtifact(g)
	out, err := yaml.Marshal(artifact)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}

// printLoaderIssues displays loader-stage issues and returns true if any errors were present.
func printLoaderIssues(issues []loader.Issue) bool {
	if len(issues) == 0 {
		return false
	}
	errs := false
	for _, is := range issues {
		printIssue("loader", is.Severity, is.Code, is.Message, "", "")
		if is.Severity == "error" {
			errs = true
		}
	}
	return errs
}

// printIssue renders a single issue consistently for loader and validator.
func printIssue(prefix, severity, code, message, gene, codon string) {
	line := fmt.Sprintf("%s-%s [%s]: %s", prefix, severity, code, message)
	if gene != "" || codon != "" {
		line += fmt.Sprintf(" (gene=%s codon=%s)", gene, codon)
	}
	fmt.Println(line)
}
