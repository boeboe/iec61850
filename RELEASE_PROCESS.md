# Release Process for v2.0.0

This guide covers the complete release process for creating v2.0.0 with libiec61850 v1.6.1.

## Overview

v2.0.0 is a **breaking change** because:
- Requires libiec61850 v1.6.1 (was v1.5.3)
- Requires mbedtls dependency for R-GOOSE and R-SMV
- All precompiled libraries must be rebuilt
- No backward compatibility with v1.5.3

## Pre-Release Checklist

### 1. Rebuild All Platform Libraries

You must rebuild precompiled libraries for **all supported platforms** to v1.6.1.

#### Option A: Automated Build (Recommended)

Use the GitHub Actions workflow to build all platforms automatically:

1. Go to your repository on GitHub
2. Click **Actions** tab
3. Select **Build libiec61850 Libraries** workflow
4. Click **Run workflow**
5. Enable **"Commit libraries to repository"**
6. Click **Run workflow**

The workflow will:
- Build all 5 platforms in parallel (~15-30 minutes)
- Automatically enable R-GOOSE, R-SMV, and SNTP
- Commit the built libraries to your repository
- Verify successful compilation

See [.github/workflows/README.md](.github/workflows/README.md) for details.

#### Option B: Manual Build

If you prefer to build locally:

**Platforms to Build:**
- ✅ **macOS ARM64** (darwin_armv8) - Build on M1/M2 Mac
- ✅ **Linux x86_64** (linux64) - Build on Linux AMD64 or cross-compile
- ✅ **Linux ARMv7l** (linux_armv7l) - Cross-compile from Linux
- ✅ **Linux ARMv8** (linux_armv8) - Cross-compile from Linux or build on ARM64
- ✅ **Windows x64** (win64) - Build on Windows with Visual Studio

**Build Steps:**

**On macOS (for darwin_armv8):**
```bash
# Install dependencies
brew install mbedtls cmake

# Run rebuild script
cd /Users/bartvanbos/Git/boeboe/iec61850
./scripts/rebuild_libraries.sh

# Verify
ls -lh libiec61850/lib/darwin_armv8/libiec61850.a
file libiec61850/lib/darwin_armv8/libiec61850.a
```

**On Linux (for linux64, linux_armv7l, linux_armv8):**
```bash
# Install dependencies
sudo apt-get update
sudo apt-get install build-essential cmake git libmbedtls-dev

# For cross-compilation, install toolchains
sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf  # ARMv7l
sudo apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu      # ARMv8

# Run rebuild script (builds all architectures if toolchains present)
cd /path/to/iec61850
./scripts/rebuild_libraries.sh

# Verify all platforms
ls -lh libiec61850/lib/*/libiec61850.a
```

**On Windows (for win64):**
```powershell
# Install dependencies (using vcpkg recommended)
vcpkg install mbedtls:x64-windows

# Clone and build libiec61850
cd C:\temp
git clone --branch v1.6.1 https://github.com/mz-automation/libiec61850.git
cd libiec61850

# Edit config/stack_config.h:
#   #define CONFIG_IEC61850_R_GOOSE 1
#   #define CONFIG_IEC61850_R_SMV 1
#   #define CONFIG_IEC61850_SNTP_CLIENT 1

mkdir build
cd build
cmake -G "Visual Studio 17 2022" -A x64 `
      -DCMAKE_BUILD_TYPE=Release `
      -DBUILD_EXAMPLES=OFF `
      -DBUILD_PYTHON_BINDINGS=OFF `
      ..
cmake --build . --config Release

# Copy to your repo (adjust path)
copy output\Release\iec61850.lib `
     C:\path\to\iec61850\libiec61850\lib\win64\libiec61850.a
```

### 2. Verify All Libraries

After building all platforms:

```bash
cd /Users/bartvanbos/Git/boeboe/iec61850

# Check all library files exist
ls -lh libiec61850/lib/*/libiec61850.a

# Verify they're recent (should show today's date)
ls -lt libiec61850/lib/*/libiec61850.a

# For the platform you're on, test version
go test -v -run TestLibraryVersion
# Expected: Using libiec61850 version: 1.6.1 (with R-GOOSE, R-SMV, SNTP enabled)
```

### 3. Update go.mod Version

Ensure `go.mod` is correct:

```bash
# Check current module path
grep ^module go.mod

# Should be:
# module github.com/boeboe/iec61850
```

### 4. Run Full Test Suite

```bash
# Build all packages
go build ./...

# Run all tests
go test -v ./...

# Test on example projects
cd test/client_rw
go test -v .
```

### 5. Update Documentation

Ensure all docs reference v2.0.0:

- [x] README.md - Version requirements
- [x] COMPATIBILITY.md - Version info
- [x] REBUILD_LIBRARIES.md - Build instructions
- [x] UPGRADE_TO_1.6.1.md - Migration guide

### 6. Create CHANGELOG

Create `CHANGELOG.md`:

```markdown
# Changelog

## [2.0.0] - 2026-02-01

### Breaking Changes
- **REQUIRED**: libiec61850 v1.6.1 (no backward compatibility with v1.5.3)
- **REQUIRED**: mbedtls dependency for R-GOOSE and R-SMV features
- Precompiled libraries rebuilt for all platforms with features enabled

### Added
- R-GOOSE (Routable GOOSE over UDP/IP) support - enabled by default
- R-SMV (Routable Sampled Values over UDP/IP) support - enabled by default
- SNTP Client (Network time synchronization) - enabled by default
- IEC 61850 Edition 2.1 support via `Edition21` constant
- Version checking via `GetLibraryVersion()` function
- Automated rebuild script with feature auto-configuration
- Comprehensive documentation (COMPATIBILITY.md, REBUILD_LIBRARIES.md, UPGRADE_TO_1.6.1.md)

### Changed
- C headers synchronized to libiec61850 v1.6.1
- MMS header structure updated (inc/mms/inc → inc/mms)
- All platform configs updated for new header paths
- Default edition changed to Edition2 in ServerConfig

### Fixed
- CGo include paths for all platforms
- Go package imports for header directories

### Migration Guide
See [UPGRADE_TO_1.6.1.md](UPGRADE_TO_1.6.1.md) for detailed migration instructions.

## [1.0.13] - Previous Release
- Based on libiec61850 v1.5.3
- Basic MMS, GOOSE L2, SV L2 support
```

## Git Tagging and Release Process

### 1. Commit All Changes

```bash
cd /Users/bartvanbos/Git/boeboe/iec61850

# Stage all changes (including binary libraries)
git add .

# Check what's being committed
git status

# Commit with descriptive message
git commit -m "Release v2.0.0: Upgrade to libiec61850 v1.6.1 with R-GOOSE, R-SMV, SNTP enabled

Breaking Changes:
- Requires libiec61850 v1.6.1 (rebuilt all platform libraries)
- Requires mbedtls for R-GOOSE and R-SMV features
- No backward compatibility with v1.5.3

New Features:
- R-GOOSE (Routable GOOSE) enabled by default
- R-SMV (Routable Sampled Values) enabled by default  
- SNTP Client enabled by default
- IEC 61850 Edition 2.1 support
- Version checking function

See UPGRADE_TO_1.6.1.md for migration guide."
```

### 2. Create Git Tag

```bash
# Create annotated tag for v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0

libiec61850 v1.6.1 with R-GOOSE, R-SMV, and SNTP enabled

Breaking Changes:
- Requires libiec61850 v1.6.1
- Requires mbedtls dependency
- All precompiled libraries rebuilt to v1.6.1

New Features:
- R-GOOSE support (Routable GOOSE over UDP/IP)
- R-SMV support (Routable Sampled Values over UDP/IP)
- SNTP Client support
- IEC 61850 Edition 2.1 support

Migration: See UPGRADE_TO_1.6.1.md"

# Verify tag
git tag -l -n9 v2.0.0
```

### 3. Push to GitHub

```bash
# Push commits
git push origin master  # or main, depending on your default branch

# Push tag
git push origin v2.0.0
```

### 4. Create GitHub Release

Go to https://github.com/boeboe/iec61850/releases/new

**Release Details:**
- **Tag**: v2.0.0 (select the tag you just pushed)
- **Release title**: v2.0.0 - libiec61850 v1.6.1 with R-GOOSE/R-SMV/SNTP
- **Description**:

```markdown
# v2.0.0 - Major Update: libiec61850 v1.6.1

## ⚠️ Breaking Changes

This is a **major version upgrade** with breaking changes:

- **REQUIRES** libiec61850 v1.6.1 (no backward compatibility with v1.5.3)
- **REQUIRES** mbedtls library for R-GOOSE and R-SMV features
- All precompiled libraries rebuilt to v1.6.1 with advanced features enabled

## ✨ New Features

### Enabled by Default:
- ✅ **R-GOOSE** - Routable GOOSE over UDP/IP
- ✅ **R-SMV** - Routable Sampled Values over UDP/IP
- ✅ **SNTP Client** - Network time synchronization

### Additional:
- IEC 61850 Edition 2.1 support (`Edition21` constant)
- Runtime version checking with `GetLibraryVersion()`
- Comprehensive documentation and migration guides

## 📦 Precompiled Libraries

All platform libraries are included and built with v1.6.1:
- macOS ARM64 (M1/M2)
- Linux x86_64
- Linux ARMv7l
- Linux ARMv8 (ARM64)
- Windows x64

## 🚀 Getting Started

```bash
# For new projects
go get github.com/boeboe/iec61850@v2.0.0

# Verify version
go run -exec echo "$(go list -m github.com/boeboe/iec61850)" .
```

## 📖 Migration Guide

Upgrading from v1.0.13? See [UPGRADE_TO_1.6.1.md](https://github.com/boeboe/iec61850/blob/v2.0.0/UPGRADE_TO_1.6.1.md)

Key changes for your external tools:
1. Update `go.mod`: `require github.com/boeboe/iec61850 v2.0.0`
2. No code changes needed (API compatible)
3. Rebuild your project: `go build`

## 📋 Full Changelog

See [CHANGELOG.md](https://github.com/boeboe/iec61850/blob/v2.0.0/CHANGELOG.md)

## 🔗 Documentation

- [README.md](https://github.com/boeboe/iec61850/blob/v2.0.0/README.md) - Overview and usage
- [COMPATIBILITY.md](https://github.com/boeboe/iec61850/blob/v2.0.0/COMPATIBILITY.md) - Feature matrix
- [REBUILD_LIBRARIES.md](https://github.com/boeboe/iec61850/blob/v2.0.0/REBUILD_LIBRARIES.md) - Build instructions

## ⚙️ Dependencies

Runtime dependencies (handled by precompiled libraries):
- libiec61850 v1.6.1 (included)
- mbedtls (statically linked in libraries)

No additional dependencies needed for users of this library.
```

**Assets**: GitHub will automatically create source code archives (zip/tar.gz)

## Updating External Tools

### For Your External CMD Tool

Update your external tool's `go.mod`:

```bash
cd /path/to/your/cmd/tool

# Update to v2.0.0
go get github.com/boeboe/iec61850@v2.0.0

# Verify
grep boeboe/iec61850 go.mod
# Should show: github.com/boeboe/iec61850 v2.0.0

# Update dependencies
go mod tidy

# Build
go build
```

### No Code Changes Required

The API is backward compatible. Your existing code should work without changes:

```go
import "github.com/boeboe/iec61850"

// All existing APIs still work
client := iec61850.NewClient(...)
server := iec61850.NewServer(...)
// etc.
```

### New Features Available

If you want to use new features:

```go
// Check version
version := iec61850.GetLibraryVersion()
fmt.Printf("Using: %s\n", version)

// Use Edition 2.1
config := iec61850.NewServerConfig()
config.Edition = iec61850.Edition21
```

## Testing the Release

### 1. Test in External Project

Create a test project:

```bash
mkdir /tmp/test-iec61850-v2
cd /tmp/test-iec61850-v2

go mod init test-project

# Add dependency
go get github.com/boeboe/iec61850@v2.0.0

# Create test file
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "github.com/boeboe/iec61850"
)

func main() {
    version := iec61850.GetLibraryVersion()
    fmt.Printf("libiec61850 version: %s\n", version)
}
EOF

# Build and run
go build
./test-project
# Expected: libiec61850 version: 1.6.1
```

### 2. Test Cross-Platform Builds

```bash
# From your test project
GOOS=linux GOARCH=amd64 go build
GOOS=linux GOARCH=arm64 go build
GOOS=linux GOARCH=arm GOARM=7 go build
GOOS=darwin GOARCH=arm64 go build
GOOS=windows GOARCH=amd64 go build
```

## Post-Release Checklist

- [ ] All platform libraries built to v1.6.1
- [ ] Version test passes (`go test -v -run TestLibraryVersion`)
- [ ] Full test suite passes (`go test -v ./...`)
- [ ] All changes committed (including binary libraries)
- [ ] v2.0.0 tag created and pushed
- [ ] GitHub release created with proper notes
- [ ] External tool tested with v2.0.0
- [ ] Cross-platform builds verified
- [ ] Documentation updated and accurate

## Binary File Policy

**YES, commit the binary files!** Here's why:

1. **User Convenience**: Users can `go get` without building C libraries
2. **Reproducible Builds**: Everyone gets the same library version
3. **Current Practice**: Your repo already does this (libiec61850/lib/*/libiec61850.a)
4. **File Sizes**: Static libraries are ~1-5MB each (acceptable for Git)

### What to Commit:
```
libiec61850/lib/darwin_armv8/libiec61850.a    # macOS ARM64
libiec61850/lib/linux64/libiec61850.a         # Linux x64
libiec61850/lib/linux_armv7l/libiec61850.a    # Linux ARM 32-bit
libiec61850/lib/linux_armv8/libiec61850.a     # Linux ARM 64-bit
libiec61850/lib/win64/libiec61850.a           # Windows x64
```

### What NOT to Commit:
- Temporary build files (*.o, *.obj)
- Build directories (build/, output/)
- Source code from libiec61850 (only headers in inc/)

## Troubleshooting

### "cannot find package" after release

User needs to update:
```bash
go get -u github.com/boeboe/iec61850@v2.0.0
go mod tidy
```

### "version mismatch" errors

Libraries weren't rebuilt properly. Re-run rebuild script for all platforms.

### Windows builds fail

Ensure Windows library is built with Visual Studio (not MinGW) for compatibility.

## Support

After release, monitor:
- GitHub Issues for user problems
- Go package documentation: https://pkg.go.dev/github.com/boeboe/iec61850@v2.0.0
- Build status on different platforms

For questions, see:
- [COMPATIBILITY.md](COMPATIBILITY.md) - Feature coverage
- [UPGRADE_TO_1.6.1.md](UPGRADE_TO_1.6.1.md) - Migration help
- GitHub Issues - Community support
