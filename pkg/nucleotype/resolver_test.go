package nucleotype

import "testing"

func TestResolvePrimitiveAlias(t *testing.T) {
	decls, err := ParseDecls(`export string = primitive`)
	if err != nil {
		t.Fatalf("parse decls: %v", err)
	}
	syms := map[string]TypeNode{"string": decls[0].Type}
	strType := NameType{Name: "string"}
	res := Resolve(strType, syms)
	if nt, ok := res.(NameType); !ok || nt.Name != "primitive" {
		t.Fatalf("expected primitive, got %#v", res)
	}
}
