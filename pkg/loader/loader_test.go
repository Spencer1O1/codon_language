package loader_test

import (
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator"
)

func TestLoadAndValidate(t *testing.T) {
	root := "../../.codon"
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := validator.Validate(g, g.TypeEnv)
	if res.HasErrors() {
		errs, warns, infos := res.Summary()
		t.Fatalf("validation failed: %d errors, %d warnings, %d infos; first: %+v", errs, warns, infos, res.Issues[0])
	}
}

func TestBuildArtifact(t *testing.T) {
	root := "../../.codon"
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := validator.Validate(g, g.TypeEnv)
	if res.HasErrors() {
		t.Fatalf("validation failed: %+v", res.Issues[0])
	}
	art := loader.BuildArtifact(g)
	if art.SchemaVersion == "" {
		t.Fatalf("schema_version missing")
	}
	if len(art.CodonSchemas) == 0 {
		t.Fatalf("codon_schemas empty")
	}
	if len(art.Nucleotypes) == 0 {
		t.Fatalf("nucleotypes empty")
	}
	if _, ok := art.Nucleotypes["string"]; !ok {
		t.Fatalf("expected nucleotype string")
	}
	if len(art.Chromosomes) == 0 {
		t.Fatalf("chromosomes empty")
	}
	// ordering deterministic: chromosomes then genes sorted
	for i := 1; i < len(art.Chromosomes); i++ {
		if art.Chromosomes[i-1].Name > art.Chromosomes[i].Name {
			t.Fatalf("chromosomes not sorted")
		}
	}
	for _, ch := range art.Chromosomes {
		if len(ch.Genes) == 0 {
			t.Fatalf("chromosome %s has no genes", ch.Name)
		}
		for i := 1; i < len(ch.Genes); i++ {
			if ch.Genes[i-1].Name > ch.Genes[i].Name {
				t.Fatalf("genes not sorted within chromosome %s", ch.Name)
			}
		}
		if ch.Genes[0].Description == "" {
			t.Fatalf("gene description missing for chromosome %s", ch.Name)
		}
	}
}
