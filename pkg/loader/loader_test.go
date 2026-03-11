package loader_test

import (
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator"
)

func TestLoadAndValidate(t *testing.T) {
	g, err := loader.LoadGenome("..")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	env, err := loader.BuildTypeEnv("..")
	if err != nil {
		t.Fatalf("type env: %v", err)
	}
	res := validator.Validate(g, env)
	if res.HasErrors() {
		errs, warns, infos := res.Summary()
		t.Fatalf("validation failed: %d errors, %d warnings, %d infos; first: %+v", errs, warns, infos, res.Issues[0])
	}
}
