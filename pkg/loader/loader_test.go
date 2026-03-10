package loader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ComposesAndOrders(t *testing.T) {
	root := filepath.Join("..", "..", "fixtures", "sample")

	cg, err := Load(root)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cg.SchemaVersion != "1.0.0" {
		t.Fatalf("schema_version = %s", cg.SchemaVersion)
	}
	if got := len(cg.Genes); got != 3 {
		t.Fatalf("expected 3 genes, got %d", got)
	}
	// Deterministic ordering: alpha.user, identity.auth, zeta.audit
	want := []struct {
		chrom string
		gene  string
	}{
		{"alpha", "user"},
		{"identity", "auth"},
		{"zeta", "audit"},
	}
	for i, w := range want {
		if cg.Genes[i].Chromosome != w.chrom || cg.Genes[i].Name != w.gene {
			t.Fatalf("gene[%d] = %s.%s, want %s.%s", i, cg.Genes[i].Chromosome, cg.Genes[i].Name, w.chrom, w.gene)
		}
	}

	user := findGene(cg, "alpha", "user")
	if user == nil {
		t.Fatalf("user gene not found")
	}
	if len(user.Entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(user.Entities))
	}
	userEntity := findEntity(user.Entities, "User")
	if userEntity == nil {
		t.Fatalf("User entity not found")
	}
	emailField := userEntity.Fields["email"]
	if !emailField.Unique || emailField.Type != "string" {
		t.Fatalf("email field parsed incorrectly: %+v", emailField)
	}
	statusField := userEntity.Fields["status"]
	if statusField.Type != "enum" || len(statusField.Values) != 2 {
		t.Fatalf("status enum parsed incorrectly: %+v", statusField)
	}

	if len(user.Capabilities) != 2 {
		t.Fatalf("expected 2 capabilities, got %d", len(user.Capabilities))
	}
	reg := findCapability(user.Capabilities, "register-user")
	if reg == nil {
		t.Fatalf("semantic capability not parsed")
	}
	createSession := findCapability(user.Capabilities, "create-session")
	if createSession == nil || createSession.Outputs["session_id"].Type != "uuid" {
		t.Fatalf("structured capability outputs not parsed: %+v", createSession)
	}

	if len(user.Relations) != 1 || user.Relations[0].Type != "many-to-one" {
		t.Fatalf("relations parsed incorrectly: %+v", user.Relations)
	}
	if len(user.References) != 1 || user.References[0].To != "identity.auth.User" {
		t.Fatalf("references parsed incorrectly: %+v", user.References)
	}
}

func TestLoad_FilenameMismatchFails(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "genome.yaml"), []byte("schema_version: 1.0.0\nproject:\n  name: temp\n"), 0o644); err != nil {
		t.Fatalf("write genome.yaml: %v", err)
	}
	chromo := filepath.Join(dir, "chromosomes", "alpha")
	if err := os.MkdirAll(chromo, 0o755); err != nil {
		t.Fatalf("mk chromo: %v", err)
	}
	genePath := filepath.Join(chromo, "user.yaml")
	geneContent := "gene:\n  name: different\n"
	if err := os.WriteFile(genePath, []byte(geneContent), 0o644); err != nil {
		t.Fatalf("write gene: %v", err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Fatalf("expected error for filename mismatch")
	}
	if !strings.Contains(err.Error(), "must match filename") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func findEntity(list []ComposedEntity, name string) *ComposedEntity {
	for i := range list {
		if list[i].Name == name {
			return &list[i]
		}
	}
	return nil
}

func findCapability(list []ComposedCapability, name string) *ComposedCapability {
	for i := range list {
		if list[i].Name == name {
			return &list[i]
		}
	}
	return nil
}

func findGene(cg *ComposedGenome, chrom, gene string) *ComposedGene {
	for i := range cg.Genes {
		if cg.Genes[i].Chromosome == chrom && cg.Genes[i].Name == gene {
			return &cg.Genes[i]
		}
	}
	return nil
}
