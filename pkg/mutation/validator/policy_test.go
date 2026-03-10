package validator

import "testing"

func TestShouldApplyStrict(t *testing.T) {
	r := &Result{}
	r.add("structural", Error, "x", "err")
	if ShouldApplyStrict(r, false) {
		t.Fatalf("should block on error")
	}

	r2 := &Result{}
	r2.add("guidance", Warning, "x", "warn")
	if !ShouldApplyStrict(r2, false) {
		t.Fatalf("warnings should pass when blockOnWarnings=false")
	}
	if ShouldApplyStrict(r2, true) {
		t.Fatalf("warnings should block when blockOnWarnings=true")
	}
}
