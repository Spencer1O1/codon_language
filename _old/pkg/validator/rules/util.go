package rules

import "fmt"

func genePath(idx int) string {
	return fmt.Sprintf("genes[%d]", idx)
}
