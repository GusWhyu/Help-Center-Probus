#!/bin/bash

# Build script untuk Layrd dengan version injection
# Jalankan di Git Bash atau WSL: ./build.sh

# Pastikan folder output bin ada
mkdir -p bin

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-s -w \
  -X main.Version=$VERSION \
  -X main.BuildTime=$BUILD_TIME \
  -X main.GitCommit=$GIT_COMMIT"

echo "Building Layrd $VERSION ($GIT_COMMIT)..."

# Build untuk Linux
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/layrd ./cmd/server
echo "  ✅ Linux   → bin/layrd"

# Build untuk Windows
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/layrd.exe ./cmd/server
echo "  ✅ Windows → bin/layrd.exe"

echo ""
echo "Build selesai! Binary tersimpan di folder bin/"
