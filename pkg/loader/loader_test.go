package loader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ComposesAndOrders(t *testing.T) {
	root := filepath.Join("..", "..", "fixtures", "shared", "valid_genome")

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
	// Deterministic ordering by chromosome then gene
	want := []struct {
		chrom string
		gene  string
	}{
		{"identity", "auth"},
		{"notifications", "notify"},
		{"tasks", "tasking"},
	}
	for i, w := range want {
		if cg.Genes[i].Chromosome != w.chrom || cg.Genes[i].Name != w.gene {
			t.Fatalf("gene[%d] = %s.%s, want %s.%s", i, cg.Genes[i].Chromosome, cg.Genes[i].Name, w.chrom, w.gene)
		}
	}

	tasks := findGene(cg, "tasks", "tasking")
	if tasks == nil {
		t.Fatalf("tasking gene not found")
	}
	if len(tasks.Entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(tasks.Entities))
	}
	taskEntity := findEntity(tasks.Entities, "Task")
	if taskEntity == nil {
		t.Fatalf("Task entity not found")
	}
	statusField := taskEntity.Fields["status"]
	if statusField.Type != "enum" || len(statusField.Values) != 3 {
		t.Fatalf("status enum parsed incorrectly: %+v", statusField)
	}

	if len(tasks.Capabilities) != 3 {
		t.Fatalf("expected 3 capabilities, got %d", len(tasks.Capabilities))
	}
	assign := findCapability(tasks.Capabilities, "assign-task")
	if assign == nil {
		t.Fatalf("assign-task capability missing")
	}
	if len(tasks.Relations) != 1 || tasks.Relations[0].Name != "parent_task" {
		t.Fatalf("relations parsed incorrectly: %+v", tasks.Relations)
	}
	if len(tasks.References) != 1 || tasks.References[0].To != "identity.auth.User" {
		t.Fatalf("references parsed incorrectly: %+v", tasks.References)
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
