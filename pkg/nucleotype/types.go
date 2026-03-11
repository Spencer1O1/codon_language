package nucleotype

import "fmt"

// TypeExpr is a placeholder for a parsed type expression string.
type TypeExpr string

// Equal compares two TypeNodes structurally.
func Equal(a, b TypeNode) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
