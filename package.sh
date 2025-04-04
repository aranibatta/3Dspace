#!/bin/bash
set -e

# Ensure fyne CLI is installed
if ! command -v $HOME/go/bin/fyne &> /dev/null
then
    echo "Installing Fyne CLI..."
    go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# Detect platform
PLATFORM="$(uname)"
echo "Detected platform: $PLATFORM"

# Create a release folder
mkdir -p releases

if [[ "$PLATFORM" == "Darwin" ]]; then
    echo "Creating macOS application..."
    $HOME/go/bin/fyne package -os darwin -icon Icon.png -name "3D Space Visualizer" -appID com.learngo.3dspace
    
    if [ -d "3D Space Visualizer.app" ]; then
        # Create DMG manually
        echo "Creating DMG..."
        hdiutil create -volname "3D Space Visualizer" -srcfolder "3D Space Visualizer.app" -ov -format UDZO "releases/3D_Space_Visualizer_macOS.dmg" 2>/dev/null || echo "DMG creation failed. macOS app available in 3D Space Visualizer.app"
        cp -r "3D Space Visualizer.app" releases/
        echo "macOS package created: releases/3D_Space_Visualizer_macOS.dmg"
    fi
elif [[ "$PLATFORM" == "Linux" ]]; then
    echo "Creating Linux application..."
    $HOME/go/bin/fyne package -os linux -icon Icon.png -name "3D Space Visualizer" -appID com.learngo.3dspace
    
    if [ -f "3D Space Visualizer" ]; then
        mv "3D Space Visualizer" releases/
        echo "Linux executable created: releases/3D Space Visualizer"
    fi
elif [[ "$PLATFORM" == "MINGW"* ]] || [[ "$PLATFORM" == "MSYS"* ]]; then
    echo "Creating Windows application..."
    $HOME/go/bin/fyne package -os windows -icon Icon.png -name "3D Space Visualizer" -appID com.learngo.3dspace
    
    if [ -f "3D Space Visualizer.exe" ]; then
        mv "3D Space Visualizer.exe" releases/
        echo "Windows executable created: releases/3D Space Visualizer.exe"
    fi
else
    echo "Unknown platform: $PLATFORM"
    exit 1
fi

echo "Packaging complete! Releases available in the 'releases' directory."