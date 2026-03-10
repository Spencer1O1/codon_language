package loader

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reChromosome = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	reGene       = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	reEntity     = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	reCapability = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	reGeneRef    = regexp.MustCompile(`^[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*$`)
	reEntityRef  = regexp.MustCompile(`^[a-z][a-z0-9_]*\.[a-z][a-z0-9_]*\.[A-Z][A-Za-z0-9]*$`)
)

var reservedWords = map[string]struct{}{
	"capabilities":   {},
	"codons":         {},
	"dependencies":   {},
	"entities":       {},
	"expression":     {},
	"gene":           {},
	"project":        {},
	"references":     {},
	"relations":      {},
	"schema":         {},
	"schema_version": {},
	"traits":         {},
	"types":          {},
}

var fieldTypes = map[string]struct{}{
	"string":    {},
	"text":      {},
	"int":       {},
	"float":     {},
	"boolean":   {},
	"uuid":      {},
	"datetime":  {},
	"json":      {},
	"enum":      {},
	"reference": {},
}

var relationTypes = map[string]struct{}{
	"one-to-one":   {},
	"one-to-many":  {},
	"many-to-one":  {},
	"many-to-many": {},
}

func isReserved(word string) bool {
	_, ok := reservedWords[word]
	return ok
}

// ValidateIdentifier enforces canonical identifier forms and reserved-word bans.
func ValidateIdentifier(kind, value string) error {
	switch kind {
	case "chromosome":
		if !reChromosome.MatchString(value) {
			return fmt.Errorf("chromosome identifier %q must match %s", value, reChromosome.String())
		}
	case "gene":
		if !reGene.MatchString(value) {
			return fmt.Errorf("gene identifier %q must match %s", value, reGene.String())
		}
	case "entity":
		if !reEntity.MatchString(value) {
			return fmt.Errorf("entity identifier %q must match %s", value, reEntity.String())
		}
	case "capability":
		if !reCapability.MatchString(value) {
			return fmt.Errorf("capability identifier %q must match %s", value, reCapability.String())
		}
	default:
		return fmt.Errorf("unknown identifier kind %q", kind)
	}

	if isReserved(value) {
		return fmt.Errorf("%s identifier %q is a reserved word", kind, value)
	}
	if len(value) == 0 || len(value) > 80 {
		return fmt.Errorf("%s identifier %q length out of bounds", kind, value)
	}
	return nil
}

// ValidateGeneReference ensures chromosome.gene format.
func ValidateGeneReference(ref string) error {
	if !reGeneRef.MatchString(ref) {
		return fmt.Errorf("dependency %q must be chromosome.gene", ref)
	}
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return fmt.Errorf("dependency %q must be chromosome.gene", ref)
	}
	if err := ValidateIdentifier("chromosome", parts[0]); err != nil {
		return err
	}
	if err := ValidateIdentifier("gene", parts[1]); err != nil {
		return err
	}
	return nil
}

// ValidateEntityReference ensures chromosome.gene.Entity format.
func ValidateEntityReference(ref string) error {
	if !reEntityRef.MatchString(ref) {
		return fmt.Errorf("entity reference %q must be chromosome.gene.Entity", ref)
	}
	parts := strings.Split(ref, ".")
	if len(parts) != 3 {
		return fmt.Errorf("entity reference %q must be chromosome.gene.Entity", ref)
	}
	if err := ValidateIdentifier("chromosome", parts[0]); err != nil {
		return err
	}
	if err := ValidateIdentifier("gene", parts[1]); err != nil {
		return err
	}
	if err := ValidateIdentifier("entity", parts[2]); err != nil {
		return err
	}
	return nil
}

func ValidateFieldType(t string) error {
	if _, ok := fieldTypes[t]; !ok {
		return fmt.Errorf("field type %q is not supported", t)
	}
	return nil
}

func ValidateRelationType(t string) error {
	if _, ok := relationTypes[t]; !ok {
		return fmt.Errorf("relation type %q is not supported", t)
	}
	return nil
}

// IsReserved reports if the word is reserved globally.
func IsReserved(word string) bool {
	return isReserved(word)
}

// GenePartFromEntityRef extracts chromosome.gene from chromosome.gene.Entity.
func GenePartFromEntityRef(ref string) (string, error) {
	if err := ValidateEntityReference(ref); err != nil {
		return "", err
	}
	parts := strings.Split(ref, ".")
	return parts[0] + "." + parts[1], nil
}

func expectKeys(name string, m map[string]any, allowed []string) error {
	allowedSet := map[string]struct{}{}
	for _, k := range allowed {
		allowedSet[k] = struct{}{}
	}
	for k := range m {
		if _, ok := allowedSet[k]; !ok {
			return fmt.Errorf("%s contains unknown field %q", name, k)
		}
	}
	return nil
}
