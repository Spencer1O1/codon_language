package core

// RuleCategory identifies a validation rule class for severity mapping.
type RuleCategory string

const (
	CategoryStructural   RuleCategory = "structural"
	CategoryIdentifiers  RuleCategory = "identifiers"
	CategoryDependencies RuleCategory = "dependencies"
	CategoryReferences   RuleCategory = "references"
	CategoryRelations    RuleCategory = "relations"
	CategoryUniqueness   RuleCategory = "uniqueness"
	CategoryTraits       RuleCategory = "traits"
	CategoryReserved     RuleCategory = "reserved_words"
)
