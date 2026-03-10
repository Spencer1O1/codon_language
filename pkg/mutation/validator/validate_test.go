package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPatchValidatorFixtures(t *testing.T) {
	cases := []struct {
		name         string
		fixture      string
		wantErrorsGE int
		wantWarnsGE  int
		state        map[string]any
	}{
		{"valid", "valid", 0, 0, nil},
		{"missing_required", "missing_required", 1, 0, nil},
		{"bad_risk", "bad_risk", 1, 0, nil},
		{"duplicate_change", "duplicate_change", 0, 1, nil}, // duplicate id is warning
		{"duplicate_ops", "duplicate_ops", 0, 1, nil},
		{"bad_op", "bad_op", 1, 0, nil},
		{"remove_no_reason", "remove_no_reason", 1, 0, nil},
		{"remove_dash", "remove_dash", 1, 0, nil},
		{"invalid_confidence", "invalid_confidence", 1, 0, nil},
		{"add_overwrite_warns", "add_overwrite", 0, 1, map[string]any{"project": map[string]any{"name": "old"}}},
		{"update_missing_errors", "update_missing", 1, 0, map[string]any{"project": map[string]any{"name": "old"}}},
		{"old_value_mismatch_errors", "old_value_mismatch", 1, 0, map[string]any{"project": map[string]any{"name": "old"}}},
		{"add_index_oob", "add_index_oob", 1, 0, map[string]any{"project": map[string]any{"tags": []any{"a"}}}},
		{"list_index_nonint", "list_index_nonint", 1, 0, map[string]any{"project": map[string]any{"tags": []any{"a"}}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "..", "fixtures", "patch", tc.fixture, "patch.yaml")
			var res *Result
			var err error
			if tc.state != nil {
				data, readErr := os.ReadFile(path)
				if readErr != nil {
					t.Fatalf("read: %v", readErr)
				}
				res, _, err = ValidateBytesWithState(data, tc.state)
			} else {
				res, _, err = ValidateFile(path)
			}
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
