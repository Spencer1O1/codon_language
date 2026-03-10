package validator

// OperationResult models execution outcome for a single operation id.
// This is a runtime companion, not the input document.
type OperationResult struct {
	ID       string
	Skipped  bool
	Severity Severity
}

// ValidateExecution checks required operations were not skipped, using the original patch doc.
func ValidateExecution(doc *PatchDoc, results []OperationResult) *Result {
	res := &Result{}
	if doc == nil {
		return res
	}
	requiredByID := map[string]bool{}
	for _, c := range doc.Changes {
		for _, op := range c.Operations {
			req := true
			if op.Required != nil {
				req = *op.Required
			}
			requiredByID[op.ID] = req
		}
	}
	resultMap := map[string]OperationResult{}
	for _, r := range results {
		resultMap[r.ID] = r
		if req, ok := requiredByID[r.ID]; ok && req && r.Skipped {
			res.add("required_op_skipped", severityByCode["required_op_skipped"], "operations["+r.ID+"]", "required operation was skipped")
		}
	}
	for id, req := range requiredByID {
		if !req {
			continue
		}
		if _, ok := resultMap[id]; !ok {
			res.add("required_op_skipped", severityByCode["required_op_skipped"], "operations["+id+"]", "required operation missing from results")
		}
	}
	return res
}
