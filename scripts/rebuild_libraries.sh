#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
LIB_DIR="$PROJECT_ROOT/libiec61850/lib"
TMP_BUILD="/tmp/libiec61850-build-$(date +%Y%m%d-%H%M%S)"

echo "=== libiec61850 Library Rebuild Script ==="
echo "Project root: $PROJECT_ROOT"
echo "Temporary build: $TMP_BUILD"
echo ""
echo "Build configuration:"
echo "  - R-GOOSE: ENABLED (requires mbedtls)"
echo "  - R-SMV: ENABLED (requires mbedtls)"
echo "  - SNTP Client: ENABLED"

# Clone libiec61850
echo ""
echo "Cloning libiec61850 v1.6.1..."
git clone --depth 1 --branch v1.6.1 \
    https://github.com/mz-automation/libiec61850.git \
    "$TMP_BUILD"

cd "$TMP_BUILD"

# Enable advanced features in stack configuration
echo ""
echo "Enabling R-GOOSE, R-SMV, and SNTP client in stack_config.h..."
sed -i.bak 's/#define CONFIG_IEC61850_R_GOOSE 0/#define CONFIG_IEC61850_R_GOOSE 1/' config/stack_config.h
sed -i.bak 's/#define CONFIG_IEC61850_R_SMV 0/#define CONFIG_IEC61850_R_SMV 1/' config/stack_config.h
sed -i.bak 's/#define CONFIG_IEC61850_SNTP_CLIENT 0/#define CONFIG_IEC61850_SNTP_CLIENT 1/' config/stack_config.h
echo "✓ Advanced features enabled"

# Function to build for a platform
build_platform() {
    local platform=$1
    local build_dir=$2
    local output_dir=$3
    shift 3
    local cmake_args=("$@")
    
    echo ""
    echo "=== Building for $platform ==="
    mkdir -p "$build_dir"
    cd "$build_dir"
    
    cmake -DCMAKE_BUILD_TYPE=Release \
          -DBUILD_EXAMPLES=OFF \
          -DBUILD_PYTHON_BINDINGS=OFF \
          "${cmake_args[@]}" \
          "$TMP_BUILD"
    
    local num_cores
    if command -v nproc &> /dev/null; then
        num_cores=$(nproc)
    elif command -v sysctl &> /dev/null; then
        num_cores=$(sysctl -n hw.ncpu)
    else
        num_cores=4
    fi
    
    make -j$num_cores
    
    mkdir -p "$output_dir"
    cp output/libiec61850.a "$output_dir/"
    echo "✓ Copied to $output_dir/libiec61850.a"
    
    cd "$TMP_BUILD"
}

# Detect current platform
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$PLATFORM" == "linux" ]] && [[ "$ARCH" == "x86_64" ]]; then
    build_platform "Linux x86_64" \
        "build-linux64" \
        "$LIB_DIR/linux64"
        
    # Check for cross-compilers
    if command -v arm-linux-gnueabihf-gcc &> /dev/null; then
        echo ""
        echo "Cross-compiling for ARMv7l..."
        
        # Create temporary toolchain file
        cat > /tmp/armv7l-toolchain.cmake << 'EOF'
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR arm)

set(CMAKE_C_COMPILER arm-linux-gnueabihf-gcc)
set(CMAKE_CXX_COMPILER arm-linux-gnueabihf-g++)

set(CMAKE_FIND_ROOT_PATH /usr/arm-linux-gnueabihf)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
EOF
        
        build_platform "Linux ARMv7l" \
            "build-armv7l" \
            "$LIB_DIR/linux_armv7l" \
            "-DCMAKE_TOOLCHAIN_FILE=/tmp/armv7l-toolchain.cmake"
    else
        echo ""
        echo "⚠ Skipping ARMv7l (cross-compiler not found)"
        echo "  Install with: sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf"
    fi
    
    if command -v aarch64-linux-gnu-gcc &> /dev/null; then
        echo ""
        echo "Cross-compiling for ARMv8..."
        
        # Create temporary toolchain file
        cat > /tmp/armv8-toolchain.cmake << 'EOF'
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR aarch64)

set(CMAKE_C_COMPILER aarch64-linux-gnu-gcc)
set(CMAKE_CXX_COMPILER aarch64-linux-gnu-g++)

set(CMAKE_FIND_ROOT_PATH /usr/aarch64-linux-gnu)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
EOF
        
        build_platform "Linux ARMv8" \
            "build-armv8" \
            "$LIB_DIR/linux_armv8" \
            "-DCMAKE_TOOLCHAIN_FILE=/tmp/armv8-toolchain.cmake"
    else
        echo ""
        echo "⚠ Skipping ARMv8 (cross-compiler not found)"
        echo "  Install with: sudo apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu"
    fi
    
elif [[ "$PLATFORM" == "darwin" ]] && [[ "$ARCH" == "arm64" ]]; then
    build_platform "macOS ARM64" \
        "build-darwin-arm64" \
        "$LIB_DIR/darwin_armv8" \
        "-DCMAKE_OSX_ARCHITECTURES=arm64"
else
    echo ""
    echo "⚠ Unsupported platform: $PLATFORM $ARCH"
    echo "Supported platforms:"
    echo "  - Linux x86_64 (with optional ARM cross-compilation)"
    echo "  - macOS ARM64 (Apple Silicon)"
    echo ""
    echo "For Windows, please build manually using Visual Studio"
    echo "See REBUILD_LIBRARIES.md for detailed instructions"
    exit 1
fi

# Cleanup
echo ""
echo "Cleaning up temporary build directory..."
rm -rf "$TMP_BUILD"
rm -f /tmp/armv7l-toolchain.cmake /tmp/armv8-toolchain.cmake

echo ""
echo "=== Build Complete ==="
echo "Libraries updated in: $LIB_DIR"
echo ""
echo "Verification:"
ls -lh "$LIB_DIR"/*/libiec61850.a 2>/dev/null || echo "  (No libraries found)"
echo ""
echo "Next steps:"
echo "1. Test version: cd $PROJECT_ROOT && go test -v -run TestLibraryVersion"
echo "2. Build project: go build ./..."
echo "3. Run full tests: go test -v ./test/..."
