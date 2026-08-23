#!/usr/bin/env bash
#
# generate.sh - Generate ogen API client from the vendored ProductBoard
# OpenAPI spec.
#
# Unlike aha-go (which hand-authors its spec because Aha publishes none),
# openapi/entities.json here is a *vendored copy* of ProductBoard's own
# official spec (see openapi/SOURCE.md for source URL, fetch date, and
# spec version — JSON has no comment syntax, so provenance lives in that
# sibling file instead of a header comment). This script only regenerates
# internal/api/ from whatever's currently vendored — it does NOT re-fetch
# the spec. Use `make update-spec` to deliberately re-fetch a newer version
# first.
#
# Usage:
#   ./generate.sh
#
# Prerequisites:
#   - Go 1.24+
#   - ogen (will be installed if missing)
#
set -euo pipefail

echo "=== productboard-go Code Generation ==="

if ! command -v ogen &> /dev/null; then
    echo "Installing ogen..."
    go install github.com/ogen-go/ogen/cmd/ogen@latest
fi

SPEC_FILE="openapi/entities.json"
if [[ ! -f "$SPEC_FILE" ]]; then
    echo "Error: $SPEC_FILE not found"
    exit 1
fi

mkdir -p internal/api

echo "Generating API client from $SPEC_FILE..."
ogen --package api --target internal/api --clean "$SPEC_FILE"

echo "Running go mod tidy..."
go mod tidy

echo "Verifying build..."
go build ./...

echo ""
echo "=== Generation complete ==="
echo ""
echo "Generated files are in internal/api/"
echo "Do not edit generated files directly."
echo ""
echo "Next steps:"
echo "  1. Update wrapper code if the spec changed"
echo "  2. Run tests: go test ./..."
echo "  3. Commit changes"
