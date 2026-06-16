#!/bin/bash
set -e

cd "$(dirname "$0")/.."
BUILD_DIR="output"
BINARY="${BUILD_DIR}/easy-flip-clock"

mkdir -p "${BUILD_DIR}"
go build -o "${BINARY}" .
exec "./${BINARY}"
