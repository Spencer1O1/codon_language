package main

import (
	"fmt"

	"github.com/Spencer1O1/codon_language/pkg/loader"
)

func runLoad(root string) error {
	g, err := loader.LoadGenome(root)
	if err != nil {
		return err
	}
	printLoaderIssues(g.Issues)
	printlnSummary(len(g.Schemas), len(g.Genes))
	return nil
}

func printlnSummary(schemas, genes int) {
	// Simple summary line for load path.
	fmt.Printf("codon schemas: %d, genes: %d\n", schemas, genes)
}
