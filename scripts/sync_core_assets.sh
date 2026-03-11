#!/usr/bin/env bash
set -euo pipefail

ROOT=${1:-.codon}
ASSETS_DIR="internal/core_assets"

SRC_FAM="$ROOT/codon_families"
SRC_TYPES="$ROOT/nucleotides/types"

DEST_FAM="$ASSETS_DIR/codon_families"
DEST_TYPES="$ASSETS_DIR/nucleotides/types"

mkdir -p "$DEST_FAM" "$DEST_TYPES"

if compgen -G "$SRC_FAM/*.codon" > /dev/null; then
  cp "$SRC_FAM"/*.codon "$DEST_FAM"/
fi

if compgen -G "$SRC_TYPES"/*.nucleotype > /dev/null; then
  cp "$SRC_TYPES"/*.nucleotype "$DEST_TYPES"/
fi

echo "Synced families -> $DEST_FAM"
echo "Synced nucleotypes -> $DEST_TYPES"
