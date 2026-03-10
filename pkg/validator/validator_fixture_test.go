package validator

import (
	"path/filepath"
	"testing"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func TestValidatorFixtures(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		wantErr  bool
		wantWarn bool
	}{
		{"valid", "valid_genome", false, false},
		{"dup_gene", "dup_gene", true, false},
		{"missing_dep", "missing_dep", true, false},
		{"reserved_word", "reserved_word", false, true},
		{"duplicate_traits", "duplicate_traits", false, true},
		{"bad_reference_target", "bad_reference_target", true, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join("..", "..", "fixtures", "validation", tc.fixture)
			cg, err := loader.Load(root)
			if err != nil {
				if !tc.wantErr {
					t.Fatalf("load error: %v", err)
				}
				// loader errors count as errors
				return
			}
			res := Validate(cg)
			errFound := res.Err() != nil
			if errFound != tc.wantErr {
				t.Fatalf("wantErr=%v, gotErr=%v, findings=%v", tc.wantErr, errFound, res.Findings)
			}
			warnFound := false
			for _, f := range res.Findings {
				if f.Severity == core.Warning {
					warnFound = true
					break
				}
			}
			if warnFound != tc.wantWarn {
				t.Fatalf("wantWarn=%v, gotWarn=%v, findings=%v", tc.wantWarn, warnFound, res.Findings)
			}
		})
	}
}
