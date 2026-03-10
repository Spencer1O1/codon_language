package loader

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoaderFixtures(t *testing.T) {
	cases := []struct {
		name    string
		fixture string
		wantErr bool
		errSub  string
	}{
		{"manifest_missing_project", "manifest_missing_project", true, "project"},
		{"missing_chromosomes", "missing_chromosomes", true, "chromosomes"},
		{"filename_mismatch", "filename_mismatch", true, "must match filename"},
		{"unknown_manifest_field", "unknown_manifest_field", true, "unknown field"},
		{"empty_genome", "empty_genome", true, "no gene files"},
		{"bad_identifier_capability", "bad_identifier_capability", true, "capability identifier"},
		{"bad_identifier_chromosome", "bad_identifier_chromosome", true, "chromosome"},
		{"missing_schema_version", "missing_schema_version", true, "schema_version"},
		{"optional_manifest_fields", "optional_manifest_fields", false, ""},
		{"valid_genome", "valid_genome", false, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join("..", "..", "fixtures", "loader", tc.fixture)
			_, err := Load(root)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tc.errSub != "" && !strings.Contains(err.Error(), tc.errSub) {
					t.Fatalf("expected error containing %q, got %v", tc.errSub, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
