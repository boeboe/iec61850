# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-02-01

### Breaking Changes

- **REQUIRED**: libiec61850 v1.6.1 (no backward compatibility with v1.5.3)
- **REQUIRED**: mbedtls dependency for R-GOOSE and R-SMV features
- Precompiled libraries rebuilt for all platforms with advanced features enabled
- Default server edition changed from Edition1 to Edition2

### Added

- **R-GOOSE** (Routable GOOSE over UDP/IP) support - enabled by default
- **R-SMV** (Routable Sampled Values over UDP/IP) support - enabled by default
- **SNTP Client** (Network time synchronization) - enabled by default
- IEC 61850 Edition 2.1 support via `Edition21` constant in `server_config.go`
- `GetLibraryVersion()` function for runtime version checking
- `version_test.go` for version verification
- Automated rebuild script (`scripts/rebuild_libraries.sh`) with feature auto-configuration
- Comprehensive documentation:
  - `COMPATIBILITY.md` - Feature coverage matrix
  - `REBUILD_LIBRARIES.md` - Library build instructions
  - `UPGRADE_TO_1.6.1.md` - Quick migration guide
  - `RELEASE_PROCESS.md` - Release workflow
- Edition constants (`Edition1`, `Edition2`, `Edition21`) for type-safe configuration

### Changed

- C headers synchronized to libiec61850 v1.6.1 from official release
- MMS header directory structure updated: `inc/mms/inc/` → `inc/mms/`
- All platform-specific config files (`config_*.go`) updated for new header paths
- Default edition in `NewServerConfig()` changed to `Edition2`
- Updated all documentation to reflect v1.6.1 requirements
- Rebuild script now auto-enables R-GOOSE, R-SMV, and SNTP features
- Enhanced error messages in version test with upgrade instructions

### Fixed

- CGo include paths for all platforms (Linux x64/ARMv7l/ARMv8, macOS ARM64, Windows x64)
- Go package imports for header directories after rsync synchronization
- Duplicate package declarations in version-related files
- MMS header path references in all build configurations

### Documentation

- README.md: Added version requirements and build instructions
- All markdown files updated to reflect v1.6.1-only support
- Removed v1.5.3 backward compatibility references
- Added quick start guide for new users

### Migration Guide

See [UPGRADE_TO_1.6.1.md](UPGRADE_TO_1.6.1.md) for detailed migration instructions from v1.0.13.

Key migration steps:
1. Update `go.mod`: `require github.com/boeboe/iec61850 v2.0.0`
2. Run `go mod tidy`
3. Rebuild your project: `go build`
4. No code changes required (API compatible)

## [1.0.13] - Previous Releases

### Based on libiec61850 v1.5.3

- MMS client/server support
- GOOSE Layer 2 (Ethernet) support
- Sampled Values Layer 2 support
- TLS support for MMS connections
- Report Control Blocks (buffered and unbuffered)
- Data Sets (static and dynamic)
- Control operations (SPC, DPC, APC, INC)
- File services (GetFile, SetFile)
- Setting Groups
- SCL/ICD file parsing
- Static model generation from SCL
- Config file model loading (.cfg files)

### Platform Support

- Linux x86_64
- Linux ARMv7l
- Linux ARMv8 (ARM64)
- macOS ARM64 (Apple Silicon)
- Windows x64

## Release Notes Format

### Version Numbering

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version (X.0.0): Incompatible API changes or dependency updates
- **MINOR** version (0.X.0): New functionality in backward-compatible manner
- **PATCH** version (0.0.X): Backward-compatible bug fixes

### Release Types

- **Breaking Changes**: Changes that require user action to upgrade
- **Added**: New features or functionality
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Now removed features
- **Fixed**: Bug fixes
- **Security**: Security vulnerability fixes

## Links

- [Repository](https://github.com/boeboe/iec61850)
- [Issue Tracker](https://github.com/boeboe/iec61850/issues)
- [Releases](https://github.com/boeboe/iec61850/releases)
- [libiec61850 Upstream](https://github.com/mz-automation/libiec61850)

[2.0.0]: https://github.com/boeboe/iec61850/releases/tag/v2.0.0
[1.0.13]: https://github.com/boeboe/iec61850/releases/tag/v1.0.13
