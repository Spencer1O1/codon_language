package nucleotype

// TypeNode is the AST interface for the codon type language.
type TypeNode interface{}

// NameType represents a bare identifier type.
type NameType struct {
	Name string
}

// LiteralType represents a string literal member (e.g., "one").
type LiteralType struct {
	Value string
}

// GenericType represents Foo<...>.
type GenericType struct {
	Name string
	Args []TypeNode
}

// OptionalType represents T?.
type OptionalType struct {
	Base TypeNode
}

// ListType represents T[].
type ListType struct {
	Base TypeNode
}

// UnionType represents T1 | T2 | ...
type UnionType struct {
	Options []TypeNode
}

// ObjectType represents { field: Type, ... }.
type ObjectType struct {
	Fields []Field
}

// Field is a named member in an ObjectType.
type Field struct {
	Name string
	Type TypeNode
}
