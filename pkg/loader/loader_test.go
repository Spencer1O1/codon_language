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
	if len(art.Genes) == 0 {
		t.Fatalf("genes empty")
	}
	// ordering deterministic: genes sorted by chromosome then gene
	for i := 1; i < len(art.Genes); i++ {
		a, b := art.Genes[i-1], art.Genes[i]
		if a.Chromosome > b.Chromosome || (a.Chromosome == b.Chromosome && a.Gene > b.Gene) {
			t.Fatalf("genes not sorted: %v then %v", a, b)
		}
	}
}
