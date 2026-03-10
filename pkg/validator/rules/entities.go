package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

// Validate entities: uniqueness and field constraints.
func init() {
	core.Register(checkEntities)
}

func checkEntities(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		index := map[string]int{}
		for ei, e := range g.Entities {
			if prev, exists := index[e.Name]; exists {
				res.Add(genePath(gi), fmt.Sprintf("duplicate entity %q (also at entities[%d])", e.Name, prev))
				continue
			}
			index[e.Name] = ei
			if err := validateEntityFields(e); err != nil {
				res.Add(fmt.Sprintf("%s.entities[%d]", genePath(gi), ei), err.Error())
			}
		}
	}
}

func validateEntityFields(e loader.ComposedEntity) error {
	for name, field := range e.Fields {
		if field.Type == "" {
			return fmt.Errorf("field %q missing type", name)
		}
		if field.Type == "enum" && len(field.Values) == 0 {
			return fmt.Errorf("field %q enum must define values", name)
		}
		if field.Type != "enum" && len(field.Values) > 0 {
			return fmt.Errorf("field %q values allowed only for enum", name)
		}
		if field.Type == "reference" && field.Reference == "" {
			return fmt.Errorf("field %q reference must be set for reference type", name)
		}
		if field.Type != "reference" && field.Reference != "" {
			return fmt.Errorf("field %q reference only valid for reference type", name)
		}
	}
	return nil
}
