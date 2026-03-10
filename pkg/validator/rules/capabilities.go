package rules

import (
	"fmt"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkCapabilities)
}

func checkCapabilities(genome *loader.ComposedGenome, res *core.Result) {
	for gi, g := range genome.Genes {
		index := map[string]int{}
		for ci, c := range g.Capabilities {
			if prev, exists := index[c.Name]; exists {
				res.Add(genePath(gi), fmt.Sprintf("duplicate capability %q (also at capabilities[%d])", c.Name, prev))
				continue
			}
			index[c.Name] = ci
			if err := validateCapability(c); err != nil {
				res.Add(fmt.Sprintf("%s.capabilities[%d]", genePath(gi), ci), err.Error())
			}
		}
	}
}

func validateCapability(c loader.ComposedCapability) error {
	for name, field := range c.Inputs {
		if field.Type == "" {
			return fmt.Errorf("input %q missing type", name)
		}
	}
	for name, field := range c.Outputs {
		if field.Type == "" {
			return fmt.Errorf("output %q missing type", name)
		}
	}
	return nil
}
