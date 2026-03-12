package core

import (
	"github.com/Spencer1O1/codon_language/pkg/loader"
	nt "github.com/Spencer1O1/codon_language/pkg/nucleotype"
)

// Rule is a validation function.
type Rule func(g *loader.Genome, env map[string]nt.TypeNode, res *Result)

var (
	registry    []Rule
	grouped     = map[string][]Rule{}
	groupOrders []string
)

// Register adds a rule to the default "language" group (backward compatible).
func Register(r Rule) {
	RegisterWithGroup("language", r)
}

// RegisterWithGroup adds a rule to a named group.
func RegisterWithGroup(group string, r Rule) {
	registry = append(registry, r)
	grouped[group] = append(grouped[group], r)
	// preserve insertion order of first appearances
	for _, g := range groupOrders {
		if g == group {
			return
		}
	}
	groupOrders = append(groupOrders, group)
}

// All returns all registered rules in registration order.
func All() []Rule {
	return registry
}

// AllByGroup returns the rules for a given group in registration order.
func AllByGroup(group string) []Rule {
	return grouped[group]
}

// Groups returns the groups in first-used order.
func Groups() []string {
	return groupOrders
}
