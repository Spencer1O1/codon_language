package validator

import (
	"path/filepath"
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Verifies rule-to-severity mapping and Result.Summary aggregation.
func TestSeverityMappingAndSummary(t *testing.T) {
	cases := []struct {
		name       string
		fixture    string
		wantErrors int
		wantWarns  int
	}{
		// missing_dep triggers two errors: missing dependency and target gene absent.
		{"missing_dep_is_error", "missing_dep", 2, 0},
		{"reserved_word_is_warning", "reserved_word", 0, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join("..", "..", "fixtures", "validation", tc.fixture)
			cg, err := loader.Load(root)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			res := Validate(cg)
			s := res.Summary()
			if s.Errors != tc.wantErrors || s.Warnings != tc.wantWarns {
				t.Fatalf("summary=%+v, want errors=%d warnings=%d", s, tc.wantErrors, tc.wantWarns)
			}
			// Err() only when errors present.
			if tc.wantErrors == 0 && res.Err() != nil {
				t.Fatalf("expected no Err(); got %v", res.Err())
			}
			if tc.wantErrors > 0 && res.Err() == nil {
				t.Fatalf("expected Err(); got nil")
			}
			// Ensure findings have expected severities.
			for _, f := range res.Findings {
				if tc.wantErrors == 0 && f.Severity == core.Error {
					t.Fatalf("unexpected error finding: %+v", f)
				}
				if tc.wantWarns == 0 && f.Severity == core.Warning {
					t.Fatalf("unexpected warning finding: %+v", f)
				}
			}
		})
	}
}
