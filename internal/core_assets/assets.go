package core_assets

import "embed"

// Embedded defaults for core codon schemas and nucleotypes.

//go:embed codon_schemas/*.yaml
var Families embed.FS

//go:embed nucleotides/types/*.nucleotype
var Nucleotypes embed.FS
