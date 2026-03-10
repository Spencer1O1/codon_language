package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func main() {
	root := flag.String("root", ".", "path to genome root")
	blockWarnings := flag.Bool("block-warnings", false, "treat validation warnings as blocking")
	flag.Parse()

	// Support optional positional root.
	if flag.NArg() > 0 && *root == "." {
		*root = flag.Arg(0)
	}

	cg, err := loader.Load(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		os.Exit(1)
	}

	res := validator.Validate(cg)
	block := false
	for _, f := range res.Findings {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", f.Severity, f.Path, f.Message)
		if f.Severity == core.Warning && *blockWarnings {
			block = true
		}
	}
	if err := res.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if block {
		fmt.Fprintf(os.Stderr, "warnings blocked validation\n")
		os.Exit(2)
	}
	fmt.Println("genome is valid")
}
