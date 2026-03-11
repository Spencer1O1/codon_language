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
	fmt.Printf("families: %d, genes: %d\n", len(g.Families), len(g.Genes))
	return nil
}

func runValidate(root string) error {
	g, err := loader.LoadGenome(root)
	if err != nil {
		return err
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
