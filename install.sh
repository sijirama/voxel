#!/bin/bash

# Check if script is run with sudo
if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo"
    exit 1
fi

# Set variables
APP_NAME="voxel"
BUILD_DIR="./build/bin"
INSTALL_DIR="/usr/local/bin"
ICON_SOURCE="./voxel.icns"  # Assuming the icon is in the same directory as the script
ICON_DIR="/usr/local/share/icons"
DESKTOP_FILE="/usr/share/applications/${APP_NAME}.desktop"

# Ensure the installation directories exist
mkdir -p "$INSTALL_DIR"
mkdir -p "$ICON_DIR"

# Build the application
echo "Building the Voxel application..."
wails build

# Check if the build was successful
if [ ! -f "$BUILD_DIR/$APP_NAME" ]; then
    echo "Build failed or executable not found. Please check your build process."
    exit 1
fi

# Make the application executable
chmod +x "$BUILD_DIR/$APP_NAME"

# Move the application to the installation directory
echo "Installing Voxel to $INSTALL_DIR..."
mv "$BUILD_DIR/$APP_NAME" "$INSTALL_DIR/"

# Copy the icon
echo "Copying icon..."
if [ -f "$ICON_SOURCE" ]; then
    cp "$ICON_SOURCE" "$ICON_DIR/${APP_NAME}.png"
else
    echo "Warning: Icon file not found. Desktop entry will be created without an icon."
fi

# Create the desktop entry
echo "Creating desktop entry..."
cat > "$DESKTOP_FILE" << EOL
[Desktop Entry]
Type=Application
Name=Voxel
Exec=$INSTALL_DIR/$APP_NAME
Icon=$ICON_DIR/${APP_NAME}.png
Categories=Utility;
EOL

# Update the desktop database
echo "Updating desktop database..."
update-desktop-database "/usr/share/applications"

echo "Installation complete. You can now run Voxel from your application launcher or by typing 'voxel' in the terminal."
