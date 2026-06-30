#!/bin/sh
set -e

REPO="duong6003/ssh-wizard"
BINARY="ssh-wizard"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# Detect arch
ARCH=$(uname -m)
case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  arm64 | aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH" && exit 1 ;;
esac

# Get latest version
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Could not fetch latest version" && exit 1
fi

FILENAME="${BINARY}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "Installing $BINARY $VERSION ($OS/$ARCH)..."

TMP=$(mktemp -d)
curl -fsSL "$URL" -o "$TMP/$FILENAME"
tar -xzf "$TMP/$FILENAME" -C "$TMP"
rm "$TMP/$FILENAME"

DEST="$INSTALL_DIR/$BINARY"
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/$BINARY" "$DEST"
else
  sudo mv "$TMP/$BINARY" "$DEST"
fi
chmod +x "$DEST"
rm -rf "$TMP"

echo "Installed to $DEST"
echo "Run: $BINARY"
