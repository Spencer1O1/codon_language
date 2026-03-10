package validator

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	return ValidateBytesWithState(data, nil)
}

// ValidateBytesWithState validates a patch document and optionally checks it against provided current state.
// state may be nil; when present it should be a map representing the target document before applying operations.
func ValidateBytesWithState(data []byte, state map[string]any) (*Result, *PatchDoc, error) {
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
	validateDocument(&doc, top, res, state)
	return res, &doc, nil
}

// ValidateFile reads a patch file and validates it.
func ValidateFile(path string) (*Result, *PatchDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read patch: %w", err)
	}
	return ValidateBytesWithState(data, nil)
}

func validateDocument(p *PatchDoc, top map[string]any, res *Result, state map[string]any) {
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
	checkAllowedKeys(res, "changes", top, []string{"summary", "rationale", "risk", "evidence", "warnings", "suggestions", "confidence", "changes"})
	seenChanges := map[string]int{}
	seenOpGlobal := map[string]string{} // opID -> change path
	for ci, c := range p.Changes {
		path := fmt.Sprintf("changes[%d]", ci)
		checkAllowedKeys(res, path, rawMap(top, "changes", ci), []string{"id", "target", "operations"})
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
			validateOperation(op, path+fmt.Sprintf(".operations[%d]", oi), res, seenOps, seenOpGlobal, path, state)
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

func validateOperation(op Operation, path string, res *Result, seenOps map[string]int, seenOpGlobal map[string]string, changePath string, state map[string]any) {
	if op.ID == "" {
		res.add("structural", severityByCode["structural"], path+".id", "operation id is required")
	} else if prev, ok := seenOps[op.ID]; ok {
		res.add("duplicate_op_id", severityByCode["duplicate_op_id"], path+".id", fmt.Sprintf("duplicate operation id %q (also at operations[%d])", op.ID, prev))
	} else if prevChange, ok := seenOpGlobal[op.ID]; ok {
		res.add("duplicate_op_id", severityByCode["duplicate_op_id"], path+".id", fmt.Sprintf("duplicate operation id %q (also in %s)", op.ID, prevChange))
	} else {
		seenOps[op.ID] = len(seenOps)
		seenOpGlobal[op.ID] = changePath
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

	// State-based checks (best effort, only when state provided).
	if state != nil {
		exists, current := resolvePath(state, op.Path)
		switch op.Op {
		case "add":
			if exists {
				res.add("add_overwrite", severityByCode["add_overwrite"], path+".path", "add targeting existing value")
			} else {
				checkAddIndexBounds(op, path+".path", res, state)
			}
		case "update", "remove":
			if !exists {
				res.add("structural", severityByCode["structural"], path+".path", "update/remove must target existing value")
			}
		}
		if op.OldValue != nil && exists {
			if !valuesEqual(op.OldValue, current) {
				res.add("old_value_mismatch", severityByCode["old_value_mismatch"], path+".old_value", "old_value does not match current value")
			}
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

// rawMap extracts map from a slice in a raw top-level map; helpers to detect unknown keys.
func rawMap(top map[string]any, key string, index int) map[string]any {
	rawList, ok := top[key].([]any)
	if !ok || index >= len(rawList) {
		return map[string]any{}
	}
	if m, ok := rawList[index].(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func isList(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}

// resolvePath navigates state using patch-style paths; returns existence and value when found.
func resolvePath(state map[string]any, path string) (bool, any) {
	if state == nil || path == "" || path == "/" {
		return false, nil
	}
	segs := strings.Split(path, "/")[1:]
	var cur any = state
	for i, seg := range segs {
		switch node := cur.(type) {
		case map[string]any:
			val, ok := node[seg]
			if !ok {
				return false, nil
			}
			cur = val
		case []any:
			if seg == "-" {
				return false, nil
			}
			idx, err := strconv.Atoi(seg)
			if err != nil || idx < 0 || idx >= len(node) {
				return false, nil
			}
			cur = node[idx]
		default:
			return false, nil
		}
		if i == len(segs)-1 {
			return true, cur
		}
	}
	return true, cur
}

func valuesEqual(a, b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// checkAddIndexBounds ensures add to list with numeric index is in bounds (<= len).
func checkAddIndexBounds(op Operation, path string, res *Result, state map[string]any) {
	segs := strings.Split(op.Path, "/")
	if len(segs) < 2 {
		return
	}
	last := segs[len(segs)-1]
	if last == "-" {
		return
	}
	idx, err := strconv.Atoi(last)
	if err != nil || idx < 0 {
		return
	}
	// resolve parent list
	parentPath := strings.Join(segs[:len(segs)-1], "/")
	exists, parent := resolvePath(state, parentPath)
	if !exists {
		return
	}
	list, ok := parent.([]any)
	if !ok {
		return
	}
	if idx > len(list) {
		res.add("path", severityByCode["path"], path, "add list index out of bounds")
	}
}
