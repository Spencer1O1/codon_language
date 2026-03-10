package nucleotype

import "testing"

func TestParseDecls(t *testing.T) {
	src := `
# comment
export field_key = Regex<"^[a-z]+$">
field<T> = Map<field_key, {
  type: T,
  default: Optional<T>,
}>
`
	decls, err := ParseDecls(src)
	if err != nil {
		t.Fatalf("parse decls: %v", err)
	}
	if len(decls) != 2 {
		t.Fatalf("expected 2 decls, got %d", len(decls))
	}
	if !decls[0].Export || decls[0].Name != "field_key" {
		t.Fatalf("unexpected first decl: %+v", decls[0])
	}
}
