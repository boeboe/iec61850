# Rebuilding libiec61850 v1.6.1 Libraries

This guide explains how to rebuild the precompiled static libraries for libiec61850 v1.6.1 with advanced features enabled.

## Build Configuration

This project requires libiec61850 v1.6.1 built with the following features **enabled by default**:
- ✅ **R-GOOSE** (Routable GOOSE over UDP/IP)
- ✅ **R-SMV** (Routable Sampled Values over UDP/IP)
- ✅ **SNTP Client** (Network time synchronization)

All three features require mbedtls to be installed on your build system.

## Prerequisites

### Common Requirements
- CMake 3.10 or higher
- Git
- C compiler (GCC, Clang, or MSVC)
- **mbedtls** (required for R-GOOSE and R-SMV)

### Platform-Specific Requirements

#### Linux
```bash
sudo apt-get update
sudo apt-get install build-essential cmake git libmbedtls-dev
```

#### macOS
```bash
brew install cmake mbedtls
# Xcode Command Line Tools
xcode-select --install
```

#### Windows
- Visual Studio 2019+ with C++ Desktop Development workload
- CMake (installed via Visual Studio or standalone)

## Quick Start

Use the automated rebuild script (automatically enables R-GOOSE, R-SMV, and SNTP):

```bash
cd /Users/bartvanbos/Git/boeboe/iec61850

# Ensure mbedtls is installed first!
# macOS: brew install mbedtls
# Linux: sudo apt-get install libmbedtls-dev

./scripts/rebuild_libraries.sh
```

## Manual Rebuild Process

### Step 1: Clone Official libiec61850

```bash
cd /tmp
git clone https://github.com/mz-automation/libiec61850.git
cd libiec61850
git checkout v1.6.1
```

### Step 2: Configure Stack Features

The rebuild script automatically enables:
- `CONFIG_IEC61850_R_GOOSE 1`
- `CONFIG_IEC61850_R_SMV 1`
- `CONFIG_IEC61850_SNTP_CLIENT 1`

If building manually, edit `config/stack_config.h`:

```c
/* Enable R-GOOSE (requires mbedtls) */
#define CONFIG_IEC61850_R_GOOSE 1

/* Enable R-SMV (requires mbedtls) */
#define CONFIG_IEC61850_R_SMV 1

/* Enable SNTP client */
#define CONFIG_IEC61850_SNTP_CLIENT 1
```

**Important**: Ensure mbedtls is installed before building:

```bash
# macOS
brew install mbedtls

# Ubuntu/Debian
sudo apt-get install libmbedtls-dev

# RHEL/CentOS  
sudo yum install mbedtls-devel
```

### Step 3: Build for Each Platform

#### Linux x86_64

```bash
cd /tmp/libiec61850
mkdir build-linux64 && cd build-linux64

cmake -DCMAKE_BUILD_TYPE=Release \
      -DBUILD_EXAMPLES=OFF \
      -DBUILD_PYTHON_BINDINGS=OFF \
      ..

make -j$(nproc)

# Copy to Go bindings
cp output/libiec61850.a \
   /Users/bartvanbos/Git/boeboe/iec61850/libiec61850/lib/linux64/
```

#### Linux ARMv7l (Cross-compilation)

Create toolchain file `cmake/armv7l-toolchain.cmake`:

```cmake
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR arm)

set(CMAKE_C_COMPILER arm-linux-gnueabihf-gcc)
set(CMAKE_CXX_COMPILER arm-linux-gnueabihf-g++)

set(CMAKE_FIND_ROOT_PATH /usr/arm-linux-gnueabihf)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
```

Install cross-compiler:
```bash
# Ubuntu
sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
```

Build:
```bash
cd /tmp/libiec61850
mkdir build-armv7l && cd build-armv7l

cmake -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_TOOLCHAIN_FILE=../cmake/armv7l-toolchain.cmake \
      -DBUILD_EXAMPLES=OFF \
      -DBUILD_PYTHON_BINDINGS=OFF \
      ..

make -j$(nproc)

cp output/libiec61850.a \
   /Users/bartvanbos/Git/boeboe/iec61850/libiec61850/lib/linux_armv7l/
```

#### Linux ARMv8 (ARM64) (Cross-compilation)

Create toolchain file `cmake/armv8-toolchain.cmake`:

```cmake
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR aarch64)

set(CMAKE_C_COMPILER aarch64-linux-gnu-gcc)
set(CMAKE_CXX_COMPILER aarch64-linux-gnu-g++)

set(CMAKE_FIND_ROOT_PATH /usr/aarch64-linux-gnu)
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
```

Install cross-compiler:
```bash
# Ubuntu
sudo apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

Build:
```bash
cd /tmp/libiec61850
mkdir build-armv8 && cd build-armv8

cmake -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_TOOLCHAIN_FILE=../cmake/armv8-toolchain.cmake \
      -DBUILD_EXAMPLES=OFF \
      -DBUILD_PYTHON_BINDINGS=OFF \
      ..

make -j$(nproc)

cp output/libiec61850.a \
   /Users/bartvanbos/Git/boeboe/iec61850/libiec61850/lib/linux_armv8/
```

#### macOS ARM64 (M1/M2)

Must build **natively on Apple Silicon Mac**:

```bash
cd /tmp/libiec61850
mkdir build-darwin-arm64 && cd build-darwin-arm64

cmake -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_OSX_ARCHITECTURES=arm64 \
      -DBUILD_EXAMPLES=OFF \
      -DBUILD_PYTHON_BINDINGS=OFF \
      ..

make -j$(sysctl -n hw.ncpu)

cp output/libiec61850.a \
   /Users/bartvanbos/Git/boeboe/iec61850/libiec61850/lib/darwin_armv8/
```

#### Windows x86_64

Must build on Windows machine with Visual Studio:

```powershell
cd C:\temp
git clone https://github.com/mz-automation/libiec61850.git
cd libiec61850
git checkout v1.6.1

mkdir build-win64
cd build-win64

cmake -G "Visual Studio 16 2019" -A x64 `
      -DCMAKE_BUILD_TYPE=Release `
      -DBUILD_EXAMPLES=OFF `
      -DBUILD_PYTHON_BINDINGS=OFF `
      ..

cmake --build . --config Release

# Copy to Go bindings (adjust path)
copy output\Release\iec61850.lib `
     C:\path\to\iec61850\libiec61850\lib\win64\libiec61850.a
```

## Automated Multi-Platform Build Script

Create `scripts/rebuild_libraries.sh`:

```bash
#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
LIB_DIR="$PROJECT_ROOT/libiec61850/lib"
TMP_BUILD="/tmp/libiec61850-build-$(date +%Y%m%d-%H%M%S)"

echo "=== libiec61850 Library Rebuild Script ==="
echo "Project root: $PROJECT_ROOT"
echo "Temporary build: $TMP_BUILD"

# Clone libiec61850
echo ""
echo "Cloning libiec61850 v1.6.1..."
git clone --depth 1 --branch v1.6.1 \
    https://github.com/mz-automation/libiec61850.git \
    "$TMP_BUILD"

cd "$TMP_BUILD"

# Function to build for a platform
build_platform() {
    local platform=$1
    local build_dir=$2
    local output_dir=$3
    local cmake_args=("${@:4}")
    
    echo ""
    echo "=== Building for $platform ==="
    mkdir -p "$build_dir"
    cd "$build_dir"
    
    cmake -DCMAKE_BUILD_TYPE=Release \
          -DBUILD_EXAMPLES=OFF \
          -DBUILD_PYTHON_BINDINGS=OFF \
          "${cmake_args[@]}" \
          "$TMP_BUILD"
    
    make -j$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)
    
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
        echo "Cross-compiling for ARMv7l..."
        build_platform "Linux ARMv7l" \
            "build-armv7l" \
            "$LIB_DIR/linux_armv7l" \
            "-DCMAKE_TOOLCHAIN_FILE=$SCRIPT_DIR/toolchains/armv7l-toolchain.cmake"
    else
        echo "⚠ Skipping ARMv7l (cross-compiler not found)"
    fi
    
    if command -v aarch64-linux-gnu-gcc &> /dev/null; then
        echo "Cross-compiling for ARMv8..."
        build_platform "Linux ARMv8" \
            "build-armv8" \
            "$LIB_DIR/linux_armv8" \
            "-DCMAKE_TOOLCHAIN_FILE=$SCRIPT_DIR/toolchains/armv8-toolchain.cmake"
    else
        echo "⚠ Skipping ARMv8 (cross-compiler not found)"
    fi
    
elif [[ "$PLATFORM" == "darwin" ]] && [[ "$ARCH" == "arm64" ]]; then
    build_platform "macOS ARM64" \
        "build-darwin-arm64" \
        "$LIB_DIR/darwin_armv8" \
        "-DCMAKE_OSX_ARCHITECTURES=arm64"
else
    echo "⚠ Unsupported platform: $PLATFORM $ARCH"
    echo "Please build manually or on appropriate platform"
fi

# Cleanup
echo ""
echo "Cleaning up temporary build directory..."
rm -rf "$TMP_BUILD"

echo ""
echo "=== Build Complete ==="
echo "Libraries updated in: $LIB_DIR"
echo ""
echo "Next steps:"
echo "1. Run: cd $PROJECT_ROOT && go test -v -run TestLibraryVersion"
echo "2. Run: go build ./..."
echo "3. Run full test suite: go test -v ./test/..."
```

Make it executable:
```bash
chmod +x scripts/rebuild_libraries.sh
```

## Verification

After rebuilding, verify the libraries and features:

```bash
cd /Users/bartvanbos/Git/boeboe/iec61850

# Check library files
ls -lh libiec61850/lib/*/libiec61850.a

# Verify library symbols (Linux) - check for R-GOOSE symbols
nm -D libiec61850/lib/linux64/libiec61850.a | grep -i goose

# Test Go bindings - should report v1.6.1
go test -v -run TestLibraryVersion

# Expected output:
# Using libiec61850 version: 1.6.1 (with R-GOOSE, R-SMV, SNTP enabled)

# Build project
go build ./...
```

## Troubleshooting

### Missing mbedtls Headers

R-GOOSE and R-SMV require mbedtls. If you get compilation errors:

```bash
# Ensure mbedtls is installed
pkg-config --modversion mbedtls

# Or manually check
ls /usr/local/include/mbedtls/
```

### Cross-Compilation Failures

For ARM cross-compilation, ensure you have the full toolchain:

```bash
# ARMv7l
sudo apt-get install gcc-arm-linux-gnueabihf \
                     g++-arm-linux-gnueabihf \
                     binutils-arm-linux-gnueabihf

# ARMv8
sudo apt-get install gcc-aarch64-linux-gnu \
                     g++-aarch64-linux-gnu \
                     binutils-aarch64-linux-gnu
```

### Windows Build Issues

- Ensure Visual Studio C++ tools are installed
- Run cmake from "Developer Command Prompt for VS"
- For MinGW-w64 (alternative):
  ```bash
  cmake -G "MinGW Makefiles" -DCMAKE_BUILD_TYPE=Release ..
  mingw32-make
  ```

## CI/CD Integration

To automate library builds in CI:

```yaml
# .github/workflows/rebuild-libs.yml
name: Rebuild Libraries

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 1 * *'  # Monthly

jobs:
  build-linux64:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        run: sudo apt-get install -y cmake build-essential
      - name: Rebuild library
        run: ./scripts/rebuild_libraries.sh
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: libiec61850-linux64
          path: libiec61850/lib/linux64/libiec61850.a
  
  build-macos-arm64:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        run: brew install cmake
      - name: Rebuild library
        run: ./scripts/rebuild_libraries.sh
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: libiec61850-darwin-arm64
          path: libiec61850/lib/darwin_armv8/libiec61850.a
```

## See Also

- [COMPATIBILITY.md](COMPATIBILITY.md) - Feature coverage matrix
- [README.md](README.md) - Project documentation
- [libiec61850 Build Documentation](https://github.com/mz-automation/libiec61850#building)
