package assets

import "embed"

// Embedded defaults for core codon families and nucleotypes.

//go:embed codon_families/*.codon
var Families embed.FS

//go:embed nucleotides/types/*.nucleotype
var Nucleotypes embed.FS
