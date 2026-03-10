package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Spencer1O1/codon-language/pkg/mutation/validator"
)

func main() {
	file := flag.String("file", "", "path to patch file (YAML/JSON)")
	blockWarnings := flag.Bool("block-warnings", false, "treat warnings as blocking")
	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "patch file is required")
		os.Exit(1)
	}

	res, _, err := validator.ValidateFile(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validate: %v\n", err)
		os.Exit(1)
	}
	if !validator.ShouldApplyStrict(res, *blockWarnings) {
		fmt.Fprintf(os.Stderr, "patch blocked; findings:\n")
		for _, f := range res.Findings {
			fmt.Fprintf(os.Stderr, "- [%s] %s: %s\n", f.Severity, f.Path, f.Message)
		}
		if res.Err() != nil {
			os.Exit(1)
		}
		// warnings-only block
		os.Exit(2)
	}
	for _, f := range res.Findings {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", f.Severity, f.Path, f.Message)
	}
	fmt.Println("patch is valid")
}
