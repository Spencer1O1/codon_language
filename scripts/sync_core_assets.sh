#!/usr/bin/env bash
set -euo pipefail

ROOT=${1:-.codon}
ASSETS_DIR="internal/core_assets"

SRC_FAM="$ROOT/codon_schemas"
SRC_TYPES="$ROOT/nucleotides/types"

DEST_FAM="$ASSETS_DIR/codon_schemas"
DEST_TYPES="$ASSETS_DIR/nucleotides/types"

mkdir -p "$DEST_FAM" "$DEST_TYPES"

if compgen -G "$SRC_FAM/*.yaml" > /dev/null; then
  cp "$SRC_FAM"/*.yaml "$DEST_FAM"/
fi

if compgen -G "$SRC_TYPES"/*.nucleotype > /dev/null; then
  cp "$SRC_TYPES"/*.nucleotype "$DEST_TYPES"/
fi

echo "Synced codon schemas -> $DEST_FAM"
echo "Synced nucleotypes   -> $DEST_TYPES"
