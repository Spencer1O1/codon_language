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
	root := fixturePath("fixtures", "validator", "language", "bad_ref")
	assertErrors(t, root, "ref_target_must_exist")
}

func TestValidate_TraitConflictShape(t *testing.T) {
	root := fixturePath("fixtures", "validator", "traits", "trait_conflict")
	assertErrors(t, root, "trait_shape_conflict")
}

func TestValidate_ManifestMissingProject(t *testing.T) {
	root := fixturePath("fixtures", "validator", "manifest", "missing_project")
	assertErrors(t, root, "project_name_required")
}

func TestValidate_IdentifierBad(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "identifier_bad")
	assertErrors(t, root, "identifier_syntax")
}

func TestValidate_IdentifierReserved(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "identifier_reserved")
	assertErrors(t, root, "identifier_reserved")
}

func TestValidate_RefTargetMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "ref_target_missing")
	assertErrors(t, root, "ref_target_must_exist")
}

func TestValidate_RelationTargetMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "relation_missing_target")
	assertErrors(t, root, "relation_target_must_exist")
}

func TestValidate_RelationBadCascade(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "relation_bad_cascade")
	assertErrors(t, root, "cascade_value_allowed")
}

func TestValidate_RelationBadOwnership(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "relation_bad_ownership")
	assertErrors(t, root, "ownership_side_must_be_valid")
}

func TestValidate_CapabilityMissingEffects(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "capability_missing_effects")
	assertErrors(t, root, "effects_required")
}

func TestValidate_EntityMissingType(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "entity_missing_type")
	assertErrors(t, root, "field_type_required")
}

func TestValidate_RefTypeUsage(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "ref_type_usage")
	assertErrors(t, root, "ref_type_usage")
}

func TestValidate_ValidationRulesEmpty(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "validation_rules_empty")
	assertErrors(t, root, "rules_required")
}

func TestValidate_CodonSchemaMissingVersion(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "codon_schema_missing_version")
	assertErrors(t, root, "schema_version_required")
}

func TestValidate_OverqualifiedWarn(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "overqualified")
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
	root := fixturePath("fixtures", "validator", "traits", "trait_merge_additive")
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
	root := fixturePath("fixtures", "validator", "traits", "trait_merge_scalar")
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
	root := fixturePath("fixtures", "validator", "traits", "trait_merge_shape")
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
	root := fixturePath("fixtures", "validator", "traits", "trait_merge_list")
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

func TestExpression_Happy(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "happy")
	_, _, res := loadAndValidate(t, root)
	if res.HasErrors() {
		t.Fatalf("expected no errors, got %+v", res.Issues)
	}
}

func TestExpression_ProjectionTargetMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "missing_target")
	assertErrors(t, root, "projection_target_exists")
}

func TestExpression_SelectorsMissing(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "missing_selectors")
	assertErrors(t, root, "projection_selectors_nonempty")
}

func TestExpression_TargetMissingKindStack(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "target_missing_fields")
	assertErrors(t, root, "target_requires_kind_and_stack")
}

func TestExpression_TemplateMissingSource(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "template_missing_source")
	assertErrors(t, root, "template_source_required")
}

func TestExpression_ParseFailedTargets(t *testing.T) {
	root := fixturePath("fixtures", "validator", "expression", "parse_failed")
	g, err := loader.LoadGenome(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	found := false
	for _, is := range g.Issues {
		if is.Code == "targets_parse_failed" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected loader issue targets_parse_failed, got %+v", g.Issues)
	}
}

func TestTypeExprDeepValidation(t *testing.T) {
	root := fixturePath("fixtures", "validator", "language", "type_expr_checks")
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
