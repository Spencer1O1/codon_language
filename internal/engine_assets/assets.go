package engine_assets

import "embed"

// Engine-only assets (not emitted to consumers).

//go:embed codon_schemas/*.yaml
var CodonSchemas embed.FS

//go:embed nucleotides/types/*.nucleotype
var Nucleotypes embed.FS
