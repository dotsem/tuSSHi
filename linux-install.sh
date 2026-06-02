#!/usr/bin/env bash

set -euo pipefail

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go is not installed. Please install Go 1.26.3 or higher." >&2
    exit 1
fi

echo "Building tusshi..."
go build -o tusshi cmd/tusshi/main.go

# determine target directories depending on privilege level
if [ "$(id -u)" -eq 0 ]; then
    BIN_DIR="/usr/local/bin"
    APP_DIR="/usr/share/applications"
    ICON_DIR="/usr/share/icons/hicolor/512x512/apps"
    PIXMAP_DIR="/usr/share/pixmaps"
    echo "Installing system-wide..."
else
    BIN_DIR="${HOME}/.local/bin"
    APP_DIR="${HOME}/.local/share/applications"
    ICON_DIR="${HOME}/.local/share/icons/hicolor/512x512/apps"
    PIXMAP_DIR="${HOME}/.local/share/icons"
    echo "Installing user-local (run with sudo for system-wide installation)..."
fi

mkdir -p "$BIN_DIR" "$APP_DIR" "$ICON_DIR" "$PIXMAP_DIR"

cp -f tusshi "$BIN_DIR/"
chmod +x "$BIN_DIR/tusshi"

if [ -f "assets/tusshi.png" ]; then
    cp -f assets/tusshi.png "$ICON_DIR/tusshi.png"
    cp -f assets/tusshi.png "$PIXMAP_DIR/tusshi.png"
fi

if [ -f "tusshi.desktop" ]; then
    cp -f tusshi.desktop "$APP_DIR/"
fi

# update desktop environment database to pick up the new shortcut
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$APP_DIR" || true
fi

echo "Installation complete!"
if [ "$(id -u)" -ne 0 ]; then
    # alert user if bin dir is not in PATH
    case ":$PATH:" in
        *:"$BIN_DIR":*) ;;
        *)
            echo "Warning: $BIN_DIR is not in your PATH."
            echo "Please add it to your shell configuration (e.g. ~/.bashrc or ~/.zshrc):"
            echo "  export PATH=\"\$PATH:$BIN_DIR\""
            ;;
    esac
fi
