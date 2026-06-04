#!/bin/bash
set -e

cd "$(dirname "$0")/.."

OUTPUT_NAME="easy-flip-clock"
BUILD_DIR="output"

mkdir -p "${BUILD_DIR}"

echo "=== Building for macOS ==="
mkdir -p "${BUILD_DIR}/macos-arm64"
mkdir -p "${BUILD_DIR}/macos-amd64"
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "${BUILD_DIR}/macos-arm64/${OUTPUT_NAME}" .
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "${BUILD_DIR}/macos-amd64/${OUTPUT_NAME}" .

mkdir -p "${BUILD_DIR}/macos-universal"
lipo -create -output "${BUILD_DIR}/macos-universal/${OUTPUT_NAME}" \
    "${BUILD_DIR}/macos-arm64/${OUTPUT_NAME}" \
    "${BUILD_DIR}/macos-amd64/${OUTPUT_NAME}"

echo "=== Building for Windows ==="
mkdir -p "${BUILD_DIR}/windows-amd64"
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "${BUILD_DIR}/windows-amd64/${OUTPUT_NAME}.exe" .

# Linux 交叉编译在 macOS 上有兼容性问题，请在 Linux CI 环境或 Docker 中编译
# CI 流水线会使用 golang:1.25 Docker 镜像在 Linux 容器中编译

echo ""
echo "=== Done ==="
ls -lh "${BUILD_DIR}"/*/