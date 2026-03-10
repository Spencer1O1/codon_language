package nucleotype

// TypeExpr is a placeholder for a parsed type expression string.
type TypeExpr string

// TypeExprOrValue captures either a type expression (for schema codons) or a raw value (for data codons).
type TypeExprOrValue struct {
	TypeExpr *string     `yaml:"type_expr,omitempty" json:"type_expr,omitempty"`
	Value    interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}
