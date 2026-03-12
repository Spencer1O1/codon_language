package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	nt "github.com/Spencer1O1/codon-language/pkg/nucleotype"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(relationsRules)
}

func relationsRules(g *loader.Genome, _ map[string]nt.TypeNode, res *core.Result) {
	for _, gene := range g.Genes {
		codon, ok := gene.Codons["relations"].(map[string]any)
		if !ok {
			continue
		}
		for name, raw := range codon {
			if strings.TrimSpace(name) == "" {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_key_required", Message: "relation names must be non-empty", Gene: gene.Name, Codon: "relations"})
				continue
			}
			rel, ok := raw.(map[string]any)
			if !ok {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_shape", Message: "relation must be an object", Gene: gene.Name, Codon: "relations"})
				continue
			}
			// ownership
			if ownRaw, ok := rel["ownership"]; ok {
				own, ok := ownRaw.(string)
				if !ok || (own != "from" && own != "to") {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "ownership_side_must_be_valid", Message: fmt.Sprintf("relation %s ownership must be 'from' or 'to'", name), Gene: gene.Name, Codon: "relations"})
				}
			}
			// cascade
			if casRaw, ok := rel["cascade"]; ok {
				cas, ok := casRaw.(string)
				if !ok || (cas != "cascade" && cas != "restrict" && cas != "nullify") {
					res.Add(core.Issue{Severity: core.SeverityError, Code: "cascade_value_allowed", Message: fmt.Sprintf("relation %s cascade must be cascade|restrict|nullify", name), Gene: gene.Name, Codon: "relations"})
				}
			}

			// from/to target existence
			from, fok := rel["from"].(string)
			to, tok := rel["to"].(string)
			if fok {
			if err := validateEntityRef(from, g, gene, res); err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_target_must_exist", Message: fmt.Sprintf("relation %s from: %v", name, err), Gene: gene.Name, Codon: "relations"})
			}
			warnRelationOverqual(from, g, gene, res, name)
		}
		if tok {
			if err := validateEntityRef(to, g, gene, res); err != nil {
				res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_target_must_exist", Message: fmt.Sprintf("relation %s to: %v", name, err), Gene: gene.Name, Codon: "relations"})
			}
			warnRelationOverqual(to, g, gene, res, name)
		}
		// relation type
		if typ, ok := rel["type"].(string); !ok || (typ != "one-to-one" && typ != "one-to-many" && typ != "many-to-one" && typ != "many-to-many") {
			res.Add(core.Issue{Severity: core.SeverityError, Code: "relation_type_allowed", Message: fmt.Sprintf("relation %s type must be one-to-one|one-to-many|many-to-one|many-to-many", name), Gene: gene.Name, Codon: "relations"})
		}
	}
}
}

var relIdentRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

// validateEntityRef reuses ref-style resolution for entity names, and warns on overqualification.
func validateEntityRef(ref string, genome *loader.Genome, gene loader.Gene, res *core.Result) error {
	parts := strings.Split(ref, ".")
	if len(parts) < 1 || len(parts) > 4 {
		return fmt.Errorf("ref must have 1-4 identifiers: %s", ref)
	}
	for _, p := range parts {
		if !relIdentRe.MatchString(p) {
			return fmt.Errorf("invalid identifier %s in %s", p, ref)
		}
	}
	switch len(parts) {
	case 1:
		// same gene, implicit entities codon
		if !hasEntry(&gene, "entities", parts[0]) {
			return fmt.Errorf("entity not found: %s", ref)
		}
		return nil
	case 2:
		// same chromosome: gene.entity
		target := findGene(genome, gene.Chromosome, parts[0])
		if !hasEntry(target, "entities", parts[1]) {
			return fmt.Errorf("entity not found: %s", ref)
		}
	case 3:
		// cross chromosome: chrom.gene.entity
		target := findGene(genome, parts[0], parts[1])
		if !hasEntry(target, "entities", parts[2]) {
			return fmt.Errorf("entity not found: %s", ref)
		}
	default:
		return fmt.Errorf("ref must be entry, gene.entry, or chrom.gene.entry: %s", ref)
	}
	return nil
}

// warnRelationOverqual emits a warning when a relation target is over-qualified.
// Minimal forms:
// - same gene: entry
// - same chromosome, different gene: gene.entry
// - cross chromosome: chrom.gene.entry
func warnRelationOverqual(ref string, genome *loader.Genome, gene loader.Gene, res *core.Result, relName string) {
	parts := strings.Split(ref, ".")
	switch len(parts) {
	case 2:
		// gene.entry but same gene
		if parts[0] == gene.Name {
			res.Add(core.Issue{Severity: core.SeverityWarn, Code: "relation_overqualified", Message: fmt.Sprintf("relation %s target %s could be shortened to %s", relName, ref, parts[1]), Gene: gene.Name, Codon: "relations"})
		}
	case 3:
		// chrom.gene.entry but same chromosome
		if parts[0] == gene.Chromosome && parts[1] == gene.Name {
			res.Add(core.Issue{Severity: core.SeverityWarn, Code: "relation_overqualified", Message: fmt.Sprintf("relation %s target %s could be shortened to %s", relName, ref, parts[2]), Gene: gene.Name, Codon: "relations"})
		}
	}
}
