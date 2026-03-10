package validator

import (
	"path/filepath"
	"testing"
)

func TestPatchValidatorFixtures(t *testing.T) {
	cases := []struct {
		name         string
		fixture      string
		wantErrorsGE int
		wantWarnsGE  int
	}{
		{"valid", "valid", 0, 0},
		{"missing_required", "missing_required", 1, 0},
		{"bad_risk", "bad_risk", 1, 0},
		{"duplicate_change", "duplicate_change", 0, 1}, // duplicate id is warning
		{"bad_op", "bad_op", 1, 0},
		{"remove_no_reason", "remove_no_reason", 1, 0},
		{"remove_dash", "remove_dash", 1, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "..", "fixtures", "patch", tc.fixture, "patch.yaml")
			res, _, err := ValidateFile(path)
			if err != nil {
				t.Fatalf("validate file: %v", err)
			}
			s := res.Summary()
			if s.Errors < tc.wantErrorsGE || s.Warnings < tc.wantWarnsGE {
				t.Fatalf("summary=%+v, want errors>=%d warnings>=%d findings=%v", s, tc.wantErrorsGE, tc.wantWarnsGE, res.Findings)
			}
			if tc.wantErrorsGE == 0 && res.Err() != nil {
				t.Fatalf("Err unexpected: %v", res.Err())
			}
			if tc.wantErrorsGE > 0 && res.Err() == nil {
				t.Fatalf("Err expected")
			}
		})
	}
}
