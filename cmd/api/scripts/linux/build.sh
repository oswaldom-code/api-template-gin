#!/bin/bash

# versión
APP_VERTION=$1

# Application name
APP_NAME="app" # Replace with your actual application name

# Output directory for compiled binaries
BUILD_DIR="bin"

# Platforms to compile for
PLATFORMS=("linux" "windows")

# Architectures to compile for
ARCHITECTURES=("386" "amd64")

# Compile the application for each platform and architecture
for os in "${PLATFORMS[@]}"; do
    for arch in "${ARCHITECTURES[@]}"; do
        echo "Compiling $os/$arch..."
        GOOS=$os GOARCH=$arch go build -o "$BUILD_DIR/$APP_VERTION/$APP_NAME-$os-$arch-v$APP_VERTION"

        # If it's Windows, add the .exe extension
        if [ "$os" = "windows" ]; then
            mv "$BUILD_DIR/$APP_VERTION/$APP_NAME-$os-$arch-v$APP_VERTION" "$BUILD_DIR/$APP_VERTION/$APP_NAME-$os-$arch-v$APP_VERTION.exe"
        fi
    done
done

echo "Compilation complete."
