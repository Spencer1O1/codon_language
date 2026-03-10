package validator

import "testing"

func TestValidateExecutionRequiredSkip(t *testing.T) {
	reqTrue := true
	doc := &PatchDoc{
		Changes: []Change{
			{
				ID:     "c1",
				Target: "genome",
				Operations: []Operation{
					{ID: "op1", Op: "add", Path: "/a"},
					{ID: "op2", Op: "add", Path: "/b", Required: &reqTrue},
				},
			},
		},
	}

	results := []OperationResult{
		{ID: "op1", Skipped: false, Severity: Info},
		{ID: "op2", Skipped: true, Severity: Warning},
	}

	res := ValidateExecution(doc, results)
	if res.Err() == nil {
		t.Fatalf("expected error for skipped required op")
	}
}

func TestValidateExecutionMissingResult(t *testing.T) {
	doc := &PatchDoc{
		Changes: []Change{
			{
				ID:     "c1",
				Target: "genome",
				Operations: []Operation{
					{ID: "op1", Op: "add", Path: "/a"},
				},
			},
		},
	}
	// no results
	res := ValidateExecution(doc, nil)
	if res.Err() == nil {
		t.Fatalf("expected error for missing required op result")
	}
}
