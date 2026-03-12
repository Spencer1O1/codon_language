package validator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func TestValidate_HappyPathExample(t *testing.T) {
	root := fixturePath("fixtures", "example")
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := Validate(g, g.TypeEnv)
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
	assertErrors(t, root, "identifier_syntax")
}

func TestValidate_IdentifierReserved(t *testing.T) {
	root := fixturePath("fixtures", "validator", "identifier_reserved", ".codon")
	assertErrors(t, root, "identifier_reserved")
}

func TestValidate_RefTargetMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "ref_target_missing", ".codon")
	assertErrors(t, root, "ref_target_must_exist")
}

func TestValidate_RelationTargetMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "relation_missing_target", ".codon")
	assertErrors(t, root, "relation_target_must_exist")
}

func TestValidate_RelationBadCascade(t *testing.T) {
	root := fixturePath("fixtures", "validator", "relation_bad_cascade", ".codon")
	assertErrors(t, root, "cascade_value_allowed")
}

func TestValidate_RelationBadOwnership(t *testing.T) {
	root := fixturePath("fixtures", "validator", "relation_bad_ownership", ".codon")
	assertErrors(t, root, "ownership_side_must_be_valid")
}

func TestValidate_CapabilityMissingEffects(t *testing.T) {
	root := fixturePath("fixtures", "validator", "capability_missing_effects", ".codon")
	assertErrors(t, root, "effects_required")
}

func TestValidate_EntityMissingType(t *testing.T) {
	root := fixturePath("fixtures", "validator", "entity_missing_type", ".codon")
	assertErrors(t, root, "field_type_required")
}

func TestValidate_RefTypeUsage(t *testing.T) {
	root := fixturePath("fixtures", "validator", "ref_type_usage", ".codon")
	assertErrors(t, root, "ref_type_usage")
}

func TestValidate_ValidationRulesEmpty(t *testing.T) {
	root := fixturePath("fixtures", "validator", "validation_rules_empty", ".codon")
	assertErrors(t, root, "rules_required")
}

func TestValidate_CodonSchemaMissingVersion(t *testing.T) {
	root := fixturePath("fixtures", "validator", "codon_schema_missing_version", ".codon")
	assertErrors(t, root, "schema_version_required")
}

func TestValidate_OverqualifiedWarn(t *testing.T) {
	root := fixturePath("fixtures", "validator", "overqualified", ".codon")
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := Validate(g, g.TypeEnv)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
	_, warns, _ := res.Summary()
	if warns == 0 {
		t.Fatalf("expected warnings for overqualified refs/relations, got none")
	}
	foundRef := false
	foundRel := false
	for _, is := range res.Issues {
		if is.Code == "ref_overqualified" {
			foundRef = true
		}
		if is.Code == "relation_overqualified" {
			foundRel = true
		}
	}
	if !foundRef || !foundRel {
		t.Fatalf("expected both ref_overqualified and relation_overqualified warnings, got %+v", res.Issues)
	}
}

func TestTraitMerge_Additive(t *testing.T) {
	root := fixturePath("fixtures", "validator", "trait_merge_additive", ".codon")
	g, env, res := loadAndValidate(t, root)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
	if _, warns, _ := res.Summary(); warns > 0 {
		t.Fatalf("expected no warnings, got %+v", res.Issues)
	}
	gene := findGene(t, g, "main", "svc")
	entities := gene.Codons["entities"].(map[string]any)
	if _, ok := entities["bar"]; !ok {
		t.Fatalf("expected injected entity bar to be present after merge")
	}
	_ = env
}

func TestTraitMerge_ScalarConflictWarn(t *testing.T) {
	root := fixturePath("fixtures", "validator", "trait_merge_scalar", ".codon")
	g, env, res := loadAndValidate(t, root)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
	foundWarn := false
	for _, is := range res.Issues {
		if is.Code == "trait_conflict_authored_wins" {
			foundWarn = true
		}
	}
	if !foundWarn {
		t.Fatalf("expected trait_conflict_authored_wins warning, got %+v", res.Issues)
	}
	gene := findGene(t, g, "main", "svc")
	entities := gene.Codons["entities"].(map[string]any)
	foo := entities["foo"].(map[string]any)
	if foo["type_expr"] != "uuid" {
		t.Fatalf("authored value should win; got %v", foo["type_expr"])
	}
	_ = env
}

func TestTraitMerge_ShapeConflict(t *testing.T) {
	root := fixturePath("fixtures", "validator", "trait_merge_shape", ".codon")
	_, _, res := loadAndValidate(t, root)
	if !res.HasErrors() {
		t.Fatalf("expected shape conflict error")
	}
	found := false
	for _, is := range res.Issues {
		if is.Code == "trait_shape_conflict" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected trait_shape_conflict, got %+v", res.Issues)
	}
}

func TestTraitMerge_ListAppend(t *testing.T) {
	root := fixturePath("fixtures", "validator", "trait_merge_list", ".codon")
	g, _, res := loadAndValidate(t, root)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
	gene := findGene(t, g, "main", "svc")
	entities := gene.Codons["entities"].(map[string]any)
	foo := entities["foo"].(map[string]any)
	tags, ok := foo["tags"].([]any)
	if !ok {
		t.Fatalf("tags should be a list, got %T", foo["tags"])
	}
	if len(tags) != 3 {
		t.Fatalf("expected list append to length 3, got %d (%v)", len(tags), tags)
	}
}

func TestTypeExprDeepValidation(t *testing.T) {
	root := fixturePath("fixtures", "validator", "type_expr_checks", ".codon")
	_, _, res := loadAndValidate(t, root)
	if !res.HasErrors() {
		t.Fatalf("expected errors for regex/map/ref violations, got none")
	}
}

// helper
func assertErrors(t *testing.T, root string, substr string) {
	t.Helper()
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := Validate(g, g.TypeEnv)
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

// loadAndValidate is a small helper shared by trait merge tests.
func loadAndValidate(t *testing.T, root string) (*loader.Genome, map[string]nt.TypeNode, core.Result) {
	t.Helper()
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res := Validate(g, g.TypeEnv)
	return g, g.TypeEnv, res
}

func findGene(t *testing.T, g *loader.Genome, chrom, name string) *loader.Gene {
	t.Helper()
	for i := range g.Genes {
		if g.Genes[i].Chromosome == chrom && g.Genes[i].Name == name {
			return &g.Genes[i]
		}
	}
	t.Fatalf("gene %s/%s not found", chrom, name)
	return nil
}
