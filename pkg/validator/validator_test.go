package validator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
)

func TestValidate_SampleFixture(t *testing.T) {
	root := filepath.Join("..", "..", "fixtures", "sample")
	cg, err := loader.Load(root)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	res := Validate(cg)
	if err := res.Err(); err != nil {
		t.Fatalf("expected valid fixture, got %v", err)
	}
}

func TestValidate_DetectsUnknownDependency(t *testing.T) {
	cg := &loader.ComposedGenome{
		SchemaVersion: "1.0.0",
		Project:       loader.Project{Name: "tmp"},
		Genes: []loader.ComposedGene{
			{
				Name:         "user",
				Chromosome:   "alpha",
				Dependencies: []string{"missing.gene"},
			},
		},
	}
	res := Validate(cg)
	if err := res.Err(); err == nil {
		t.Fatalf("expected error for missing dependency")
	} else if !strings.Contains(err.Error(), "missing.gene") {
		t.Fatalf("unexpected error: %v", err)
	}
}
