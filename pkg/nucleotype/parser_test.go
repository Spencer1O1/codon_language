package nucleotype

import "testing"

func TestParseBasics(t *testing.T) {
	cases := []string{
		"string",
		"Map<entity_name, field<object>>",
		"Union<\"a\", \"b\">",
		"List<string>",
		"{ from: entity_name, to: entity_name, type: Union<\"one-to-one\", \"one-to-many\"> }",
		"field<object>[]?",
	}
	for _, c := range cases {
		if _, err := Parse(c); err != nil {
			t.Fatalf("parse %q: %v", c, err)
		}
	}
}

func TestParsePretty(t *testing.T) {
	src := "Map<capability_name, { effects: List<string>, inputs: Optional<object>, outputs: Optional<object> }>"
	ast, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := Pretty(ast)
	if got == "" {
		t.Fatalf("Pretty returned empty")
	}
}
