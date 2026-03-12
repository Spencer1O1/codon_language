#!/usr/bin/env bash
set -euo pipefail

ROOT=${1:-.codon}
CORE_ASSETS="internal/core_assets"
ENGINE_ASSETS="internal/engine_assets"

SRC_FAM="$ROOT/codon_schemas"
SRC_TYPES="$ROOT/nucleotides/types"

CORE_FAM_DEST="$CORE_ASSETS/codon_schemas"
CORE_TYPES_DEST="$CORE_ASSETS/nucleotides/types"
ENGINE_FAM_DEST="$ENGINE_ASSETS/codon_schemas"
ENGINE_TYPES_DEST="$ENGINE_ASSETS/nucleotides/types"

mkdir -p "$CORE_FAM_DEST" "$CORE_TYPES_DEST" "$ENGINE_FAM_DEST" "$ENGINE_TYPES_DEST"

# Clean previous copies
rm -f "$CORE_TYPES_DEST"/*.nucleotype "$ENGINE_TYPES_DEST"/*.nucleotype >/dev/null 2>&1 || true
rm -f "$CORE_FAM_DEST"/*.yaml "$ENGINE_FAM_DEST"/*.yaml >/dev/null 2>&1 || true

# Public codon schemas: core.yaml only
if [ -f "$SRC_FAM/core.yaml" ]; then
  cp "$SRC_FAM/core.yaml" "$CORE_FAM_DEST"/
fi

# Engine codon schemas: everything else (e.g., custom.yaml)
for f in "$SRC_FAM"/*.yaml; do
  [ "$(basename "$f")" = "core.yaml" ] && continue
  [ -f "$f" ] || continue
  cp "$f" "$ENGINE_FAM_DEST"/
done

# Public nucleotypes: primitives only
if [ -f "$SRC_TYPES/primitives.nucleotype" ]; then
  cp "$SRC_TYPES/primitives.nucleotype" "$CORE_TYPES_DEST"/
fi

# Engine nucleotypes: all non-primitives
for f in "$SRC_TYPES"/*.nucleotype; do
  [ -f "$f" ] || continue
  [ "$(basename "$f")" = "primitives.nucleotype" ] && continue
  cp "$f" "$ENGINE_TYPES_DEST"/
done

echo "Synced public codon schemas   -> $CORE_FAM_DEST"
echo "Synced engine codon schemas   -> $ENGINE_FAM_DEST"
echo "Synced public nucleotypes     -> $CORE_TYPES_DEST (primitives only)"
echo "Synced engine nucleotypes     -> $ENGINE_TYPES_DEST (non-primitives)"
