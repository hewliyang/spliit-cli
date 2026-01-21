#!/bin/bash
set -e

REPO="hewliyang/spliit-cli"
BINARY="spliit"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin*) OS="darwin" ;;
  linux*) OS="linux" ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "^location:" | sed 's/.*tag\///' | tr -d '\r\n')
if [ -z "$VERSION" ]; then
  echo "Failed to get latest version"
  exit 1
fi

echo "Installing $BINARY $VERSION ($OS/$ARCH)..."

# Set extension
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
  BINARY="spliit.exe"
fi

# Download URL
URL="https://github.com/$REPO/releases/download/$VERSION/spliit_${VERSION#v}_${OS}_${ARCH}.${EXT}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download and extract
echo "Downloading from $URL"
curl -sL "$URL" -o "$TMP_DIR/archive.$EXT"

if [ "$EXT" = "zip" ]; then
  unzip -q "$TMP_DIR/archive.$EXT" -d "$TMP_DIR"
else
  tar -xzf "$TMP_DIR/archive.$EXT" -C "$TMP_DIR"
fi

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
else
  echo "Need sudo to install to $INSTALL_DIR"
  sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/"
fi

chmod +x "$INSTALL_DIR/$BINARY"

echo "✓ Installed $BINARY to $INSTALL_DIR/$BINARY"
echo "  Run 'spliit --help' to get started"
