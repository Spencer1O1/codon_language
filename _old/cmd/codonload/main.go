package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"gopkg.in/yaml.v3"
)

func main() {
	root := flag.String("root", ".", "path to genome composition root")
	format := flag.String("format", "json", "output format: json|yaml")
	flag.Parse()

	cg, err := loader.Load(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		os.Exit(1)
	}

	switch *format {
	case "json":
		out, err := json.MarshalIndent(cg, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	case "yaml", "yml":
		out, err := yaml.Marshal(cg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal yaml: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q (use json|yaml)\n", *format)
		os.Exit(1)
	}
}
