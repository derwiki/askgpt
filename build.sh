#!/bin/bash

# Build for Windows (64-bit)
echo "Building for Windows (64-bit)"
GOOS=windows GOARCH=amd64 go build -o release/askgpt.exe main.go


# Build for Linux (64-bit)
echo "Building for Linux (64-bit)"
GOOS=linux GOARCH=amd64 go build -o release/askgpt-linux main.go
chmod +x release/askgpt-linux

# Build for macOS (64-bit arm)
echo "Building for macOS (64-bit)"
GOOS=darwin GOARCH=arm64 go build -o release/askgpt-macos-apple main.go
chmod +x release/askgpt-macos-apple

# Build for macOS (64-bit intel)
echo "Building for macOS (64-bit)"
GOOS=darwin GOARCH=amd64 go build -o release/askgpt-macos-intel main.go
chmod +x release/askgpt-macos-intel

# Build for FreeBSD (64-bit)
echo "Building for FreeBSD (64-bit)"
GOOS=freebsd GOARCH=amd64 go build -o release/askgpt-freebsd main.go
chmod +x release/askgpt-freebsd

# Build for OpenBSD (64-bit)
echo "Building for OpenBSD (64-bit)"
GOOS=openbsd GOARCH=amd64 go build -o release/askgpt-openbsd main.go
chmod +x release/askgpt-openbsd

echo "Done!"

