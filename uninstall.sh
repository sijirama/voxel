#!/bin/bash

# Check if script is run with sudo
if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo"
    exit 1
fi

# Set variables
APP_NAME="voxel"
INSTALL_DIR="/usr/local/bin"
ICON_DIR="/usr/local/share/icons"
DESKTOP_FILE="/usr/share/applications/${APP_NAME}.desktop"

# Remove the application
echo "Removing Voxel from $INSTALL_DIR..."
rm -f "$INSTALL_DIR/$APP_NAME"

# Remove the icon
echo "Removing icon..."
rm -f "$ICON_DIR/${APP_NAME}.icns"
rm -f "$ICON_DIR/${APP_NAME}.png"

# Remove the desktop entry
echo "Removing desktop entry..."
rm -f "$DESKTOP_FILE"

# Update the desktop database
echo "Updating desktop database..."
update-desktop-database "/usr/share/applications"

echo "Uninstallation complete. Voxel has been removed from your system."
