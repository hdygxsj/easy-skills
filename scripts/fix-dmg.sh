#!/bin/bash
# Script to rebuild dmg with frontend assets included

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

APP_PATH="$PROJECT_DIR/src-tauri/target/release/bundle/macos/Easy Skills.app"
DMG_PATH="$PROJECT_DIR/src-tauri/target/release/bundle/dmg"
DMG_NAME="Easy Skills_0.1.0_aarch64.dmg"
TEMP_DMG="/tmp/easy-skills-temp.dmg"
FINAL_DMG="$DMG_PATH/$DMG_NAME"

# Check if app exists
if [ ! -d "$APP_PATH" ]; then
    echo "Error: App not found at $APP_PATH"
    exit 1
fi

# Remove old dmg if exists
rm -f "$FINAL_DMG"

# Create temp directory for dmg content
TEMP_DIR=$(mktemp -d)
cp -r "$APP_PATH" "$TEMP_DIR/"

# Create new dmg
hdiutil create -volname "Easy Skills" -srcfolder "$TEMP_DIR" -ov -format UDZO "$TEMP_DMG"

# Move to final location
mkdir -p "$DMG_PATH"
mv "$TEMP_DMG" "$FINAL_DMG"

# Cleanup
rm -rf "$TEMP_DIR"

echo "Created dmg with frontend assets: $FINAL_DMG"
