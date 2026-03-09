package validation

import "github.com/Spencer1O1/codon-language/internal/domain/genome"

type Rule interface {
	Name() string
	Validate(*genome.Genome) []Finding
}
