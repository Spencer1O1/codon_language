package validator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
)

func TestValidate_HappyPathExample(t *testing.T) {
	root := fixturePath("fixtures", "example")
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	env, err := loader.BuildTypeEnv(root)
	if err != nil {
		t.Fatalf("env: %v", err)
	}
	res := Validate(g, env)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
	_, warns, _ := res.Summary()
	if warns > 0 {
		t.Fatalf("expected no warnings, got %+v", res.Issues)
	}
}

func TestValidate_BadRef(t *testing.T) {
	root := fixturePath("fixtures", "validator", "bad_ref", ".codon")
	assertErrors(t, root, "ref_target_must_exist")
}

func TestValidate_TraitConflictShape(t *testing.T) {
	root := fixturePath("fixtures", "validator", "trait_conflict", ".codon")
	assertErrors(t, root, "trait_shape_conflict")
}

func TestValidate_ManifestMissingProject(t *testing.T) {
	root := fixturePath("fixtures", "validator", "manifest_missing_project", ".codon")
	assertErrors(t, root, "project_name_required")
}

func TestValidate_IdentifierBad(t *testing.T) {
	root := fixturePath("fixtures", "validator", "identifier_bad", ".codon")
	assertErrors(t, root, "identifier_invalid")
}

// helper
func assertErrors(t *testing.T, root string, substr string) {
	t.Helper()
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	env, err := loader.BuildTypeEnv(root)
	if err != nil {
		t.Fatalf("env: %v", err)
	}
	res := Validate(g, env)
	if !res.HasErrors() {
		t.Fatalf("expected errors containing %s, got none", substr)
	}
	found := false
	for _, is := range res.Issues {
		if strings.Contains(is.Code, substr) || strings.Contains(is.Message, substr) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected error containing %q, got %+v", substr, res.Issues)
	}
}

// fixturePath builds a path relative to repo root (tests run from pkg/validator).
func fixturePath(parts ...string) string {
	all := append([]string{"..", ".."}, parts...)
	return filepath.Join(all...)
}
