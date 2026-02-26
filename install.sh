#!/bin/sh
set -e

# md-fetch installation script
# Detects OS and Architecture to download the correct binary from GitHub releases.
# Installs to $HOME/.local/bin by default.

OWNER="nathabonfim59"
REPO="md-fetch"
BINARY_NAME="md-fetch"
INSTALL_DIR="${HOME}/.local/bin"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Find latest version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not find latest version for $REPO"
    exit 1
fi

# Remove 'v' prefix for filename
VERSION_NUMBER=$(echo "$LATEST_VERSION" | sed 's/^v//')

# OS detection
OS=$(uname -s)
case "$OS" in
    Linux)  OS_NAME="Linux" ;;
    Darwin) OS_NAME="Darwin" ;;
    *)      echo "Error: OS $OS not supported"; exit 1 ;;
esac

# Arch detection
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH_NAME="x86_64" ;;
    arm64|aarch64) ARCH_NAME="arm64" ;;
    *)      echo "Error: Architecture $ARCH not supported"; exit 1 ;;
esac

FILENAME="${BINARY_NAME}_${VERSION_NUMBER}_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$LATEST_VERSION/$FILENAME"

echo "Downloading $BINARY_NAME $LATEST_VERSION for $OS_NAME $ARCH_NAME..."
curl -sSL "$DOWNLOAD_URL" -o "$FILENAME"

echo "Extracting..."
tar -xzf "$FILENAME" "$BINARY_NAME"
rm "$FILENAME"

chmod +x "$BINARY_NAME"

mv "$BINARY_NAME" "$INSTALL_DIR/"

echo "Successfully installed $BINARY_NAME to $INSTALL_DIR/"

# Check if INSTALL_DIR is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "Warning: $INSTALL_DIR is not in your PATH."
    echo "You may need to add it to your shell configuration (e.g., .bashrc or .zshrc):"
    echo "  export PATH="\$PATH:$INSTALL_DIR""
fi

echo "Done! Run '$BINARY_NAME --help' to get started."
