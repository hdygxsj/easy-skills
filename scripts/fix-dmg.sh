#!/bin/bash
# Script to rebuild dmg with frontend assets included

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CREATE_DMG="$PROJECT_DIR/src-tauri/target/release/bundle/dmg/bundle_dmg.sh"

APP_PATH="$PROJECT_DIR/src-tauri/target/release/bundle/macos/Easy Skills.app"
DMG_DIR="$PROJECT_DIR/src-tauri/target/release/bundle/dmg"
FINAL_DMG="$DMG_DIR/Easy Skills_0.1.0_aarch64.dmg"
TEMP_DIR="$PROJECT_DIR/tmp_dmg"

# Check if app exists
if [ ! -d "$APP_PATH" ]; then
    echo "Error: App not found at $APP_PATH"
    exit 1
fi

# Check if create-dmg script exists
if [ ! -f "$CREATE_DMG" ]; then
    echo "Error: create-dmg script not found at $CREATE_DMG"
    exit 1
fi

# Clean up old files
rm -rf "$TEMP_DIR" "$FINAL_DMG"
mkdir -p "$TEMP_DIR"

# Copy app to temp dir
cp -r "$APP_PATH" "$TEMP_DIR/"

# Use create-dmg with standard Mac installer style
"$CREATE_DMG" \
    --volname "Easy Skills" \
    --volicon "$APP_PATH/Contents/Resources/Easy Skills.icns" \
    --window-pos 200 120 \
    --window-size 540 400 \
    --icon-size 96 \
    --icon "Easy Skills.app" 160 205 \
    --app-drop-link 420 205 \
    --hide-extension "Easy Skills.app" \
    "$FINAL_DMG" \
    "$TEMP_DIR"

# Cleanup
rm -rf "$TEMP_DIR"

echo "Created dmg with Mac installer style: $FINAL_DMG"
