package core_assets

import "embed"

// Public assets embedded for consumers (core codon schemas + primitives only).

//go:embed codon_schemas/*.yaml
var CodonSchemas embed.FS

//go:embed nucleotides/types/primitives.nucleotype
var Nucleotypes embed.FS
