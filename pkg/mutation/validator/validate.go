package validator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"gopkg.in/yaml.v3"
)

var (
	riskLevels       = map[string]struct{}{"low": {}, "medium": {}, "high": {}}
	confidenceLevels = map[string]struct{}{"low": {}, "medium": {}, "high": {}}
	opTypes          = map[string]struct{}{"add": {}, "update": {}, "remove": {}}
)

// severity mapping per mutation.validation.severity_mapping
var severityByCode = map[string]Severity{
	"structural":          Error,
	"path":                Error,
	"op_type":             Error,
	"old_value_mismatch":  Error,
	"required_op_skipped": Error,
	"duplicate_op_id":     Warning,
	"add_overwrite":       Warning,
	"merge_attempt":       Error,
	"guidance":            Info,
}

// ValidateBytes validates a patch document given as YAML or JSON bytes.
func ValidateBytes(data []byte) (*Result, *PatchDoc, error) {
	var top map[string]any
	if err := yaml.Unmarshal(data, &top); err != nil {
		return nil, nil, fmt.Errorf("parse patch: %w", err)
	}

	res := &Result{}
	checkAllowedKeys(res, "patch", top, []string{"summary", "rationale", "risk", "evidence", "warnings", "suggestions", "confidence", "changes"})

	var doc PatchDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, nil, fmt.Errorf("decode patch: %w", err)
	}
	validateDocument(&doc, res)
	return res, &doc, nil
}

// ValidateFile reads a patch file and validates it.
func ValidateFile(path string) (*Result, *PatchDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read patch: %w", err)
	}
	return ValidateBytes(data)
}

func validateDocument(p *PatchDoc, res *Result) {
	if p.Summary == "" {
		res.add("structural", severityByCode["structural"], "summary", "summary is required")
	}
	if p.Rationale == "" {
		res.add("structural", severityByCode["structural"], "rationale", "rationale is required")
	}
	if p.Risk == "" {
		res.add("structural", severityByCode["structural"], "risk", "risk is required")
	} else if _, ok := riskLevels[p.Risk]; !ok {
		res.add("structural", severityByCode["structural"], "risk", "risk must use mutation.patch_base.risk_level")
	}
	if p.Confidence != "" {
		if _, ok := confidenceLevels[p.Confidence]; !ok {
			res.add("structural", severityByCode["structural"], "confidence", "confidence must use mutation.patch_base.confidence_level")
		}
	}
	if p.Warnings != nil && !isList(p.Warnings) {
		res.add("structural", severityByCode["structural"], "warnings", "warnings must be a list when present")
	}
	if p.Suggestions != nil && !isList(p.Suggestions) {
		res.add("structural", severityByCode["structural"], "suggestions", "suggestions must be a list when present")
	}
	if len(p.Changes) == 0 {
		res.add("structural", severityByCode["structural"], "changes", "changes are required")
	}
	seenChanges := map[string]int{}
	for ci, c := range p.Changes {
		path := fmt.Sprintf("changes[%d]", ci)
		if c.ID == "" {
			res.add("structural", severityByCode["structural"], path+".id", "change id is required")
		} else if prev, ok := seenChanges[c.ID]; ok {
			res.add("duplicate_op_id", severityByCode["duplicate_op_id"], path+".id", fmt.Sprintf("duplicate change id %q (also at changes[%d])", c.ID, prev))
		} else {
			seenChanges[c.ID] = ci
		}
		validateTarget(c.Target, path+".target", res)
		if len(c.Operations) == 0 {
			res.add("structural", severityByCode["structural"], path+".operations", "operations are required")
		}
		seenOps := map[string]int{}
		for oi, op := range c.Operations {
			validateOperation(op, path+fmt.Sprintf(".operations[%d]", oi), res, seenOps)
		}
	}
}

func validateTarget(target, path string, res *Result) {
	if target == "" {
		res.add("structural", severityByCode["structural"], path, "target is required")
		return
	}
	if target != "genome" {
		if err := loader.ValidateGeneReference(target); err != nil {
			res.add("structural", severityByCode["structural"], path, "target must be \"genome\" or chromosome.gene")
		}
	}
}

func validateOperation(op Operation, path string, res *Result, seenOps map[string]int) {
	if op.ID == "" {
		res.add("structural", severityByCode["structural"], path+".id", "operation id is required")
	} else if prev, ok := seenOps[op.ID]; ok {
		res.add("duplicate_op_id", severityByCode["duplicate_op_id"], path+".id", fmt.Sprintf("duplicate operation id %q (also at operations[%d])", op.ID, prev))
	} else {
		seenOps[op.ID] = len(seenOps)
	}

	if op.Op == "" {
		res.add("op_type", severityByCode["op_type"], path+".op", "op is required")
	} else if _, ok := opTypes[op.Op]; !ok {
		res.add("op_type", severityByCode["op_type"], path+".op", "op must be add, update, or remove")
	}

	if op.Path == "" {
		res.add("path", severityByCode["path"], path+".path", "path is required")
	} else {
		validatePath(op, path+".path", res)
	}

	// destructive operations must include reason
	if op.Op == "remove" && op.Reason == "" {
		res.add("structural", severityByCode["structural"], path+".reason", "destructive operations must include a reason")
	}

	// list removals must be deterministic: last segment cannot be "-" for remove/update
	if op.Op == "remove" || op.Op == "update" {
		if strings.HasSuffix(op.Path, "/-") {
			res.add("path", severityByCode["path"], path+".path", "remove/update must not use '-' list append segment")
		}
	}
}

var pathRe = regexp.MustCompile(`^/[\S]*$`)

func validatePath(op Operation, path string, res *Result) {
	if !pathRe.MatchString(op.Path) {
		res.add("path", severityByCode["path"], path, "path must start with '/' and contain no spaces")
		return
	}
	segments := strings.Split(op.Path, "/")[1:]
	if len(segments) == 0 {
		res.add("path", severityByCode["path"], path, "path must address a location")
		return
	}
	if op.Op != "add" && segments[len(segments)-1] == "-" {
		res.add("path", severityByCode["path"], path, "'-' list append allowed only for add operations")
	}
}

func checkAllowedKeys(res *Result, scope string, m map[string]any, allowed []string) {
	allow := map[string]struct{}{}
	for _, k := range allowed {
		allow[k] = struct{}{}
	}
	for k := range m {
		if _, ok := allow[k]; !ok {
			res.add("structural", severityByCode["structural"], scope, fmt.Sprintf("unknown field %q", k))
		}
	}
}

func isList(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}
