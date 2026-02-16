# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.6] - 2026-02-15

### Added

- **Server – timestamp attributes**
  - `UpdateTimestampAttributeValue(node *ModelNode, ts *Timestamp)` – update a UTC time (timestamp) DataAttribute with full Timestamp (time + time quality); use when both value and quality must be set (otherwise `UpdateUTCTimeAttributeValue(node, ms)` only sets the time value)
- **Timestamp – building from epoch ms and quality byte**
  - `SetTimeInMilliseconds(epochMs int64)` – set time from milliseconds since Unix epoch
  - `SetTimeQuality(timeQuality uint8)` – set IEC 61850 time quality from one byte (bit 7 = leap seconds known, bit 6 = clock failure, bit 5 = clock not synchronized, bits 0–4 = subsecond precision)
  - `SetSubsecondPrecision(precision int)` – set the number of significant bits of the fraction-of-second part (IEC 61850 time quality bits 0–4)

---

## [1.1.3] - 2026-02-15

### Added

- **GOOSE – publisher**
  - `CommParameters` – named Go struct (VlanPriority, VlanID, AppID, DstAddr) matching C `struct sCommParameters`; embedded in `GoosePublisherConf` so existing config literals remain valid
  - `NewGoosePublisherEx(conf, useVlanTag)` – create publisher with optional VLAN tag
  - `SetGoID(goID)` – set GOOSE identifier in messages
  - `PublishAndDump(dataSet, msgBuf)` – publish and copy raw payload into buffer
- **GOOSE – receiver**
  - `NewGooseReceiverEx(buffer)` – receiver using provided buffer for message handling
  - `StartThreadless()` / `StopThreadless()` – non-threaded operation with external read loop
  - `HandleMessage(buffer)` – parse GOOSE from raw Ethernet frame
- **GOOSE – subscriber**
  - `SetObserver()` – listen to any received GOOSE message
  - `NewGooseSubscriberWithDataSet(conf, dataSetValues)` – subscriber writing into pre-allocated data set (e.g. from `Client.ReadDataSetValues` + `ClientDataSet.GooseDataSetValues()`)
- **Client – GOOSE and data sets**
  - `ReadDataSetValues(dataSetReference)` – returns `*ClientDataSet`; call `Destroy()` when done
  - `ClientDataSet.GooseDataSetValues()` – returns `*GooseDataSetValues` for use with `NewGooseSubscriberWithDataSet`
  - `GetGoCBValuesAsync(goCBReference, callback)` / `SetGoCBValuesAsync(goCBReference, values, ...)` – async GoCB read/write
- **Server – GOOSE**
  - `EnableGoosePublishing()` / `DisableGoosePublishing()` – enable/disable integrated GOOSE publisher
  - `SetGooseInterfaceId(interfaceId)` – set Ethernet interface for GOOSE (e.g. `"eth0"`)

### Changed

- **Documentation**
  - GAPS.md – new “GOOSE API coverage” section (summary table, C→Go quick reference, link to GOOSE_GAPS.md)
  - GOOSE_GAPS.md – new file: GOOSE gap analysis and plan; executive summary and Part 1/2 (~95% Ethernet GOOSE coverage, structs/enums 100%)
  - STRUCTS.md – GOOSE section: CommParameters, GoosePublisherConf (embedding), ClientGooseControlBlock (opaque), ClientGooseControlBlockValues
  - FUNCTIONS.md, ENUMS.md – GOOSE-related entries added/updated

---

## [1.1.1] - 2026-02-07

### Added

- **Client – deferred connect and auth**
  - `NewClientWithoutConnect()` / `NewClientWithoutConnectWithTls()` – create client without connecting
  - `ConnectWithAuth(hostname, port, username, password)` – connect with ACSE password authentication
  - `ConnectAsync(hostname, port, callback)` – non-blocking connect with completion callback
- **Client – file and MMS**
  - `GetFileDirectoryExEntries()` – file directory as `[]MmsFileDirectoryEntryEx` (Filename, FileSize, LastModifiedTime, FileAttributes)
  - `Write()` accepts `*MmsValueRef` for pass-through values (no double-free)
- **Low-level MMS client** (`mms_connection.go`, `client_mms.go`)
  - `MmsConnection` – create, connect (sync/async), TLS, timeouts, disconnect/abort/conclude
  - ISO and MMS connection parameters: `IsoConnectionParameters`, `MmsConnectionParameters`, getters/setters
  - `SetRawMessageHandler()`, `SetConnectionLostHandler()`, `SetInformationReportHandler()`, `SetFilestoreBasepath()`
  - Full MMS client operations: read, write, get variable list, get name list, file services, journal client
- **MmsValue and type system** (`mms_value.go`, `mms_type_spec.go`, `types.go`)
  - `MmsValue` / `MmsValueRef` – structured MMS values with constructors, getters, setters
  - `MmsType`, `MmsDataAccessError`, bit-string and octet-string helpers, array/structure ops
  - `MmsTypeSpec` and type specification types for structured data
  - `CMmsValueToMmsValue()` – convert C `MmsValue*` to Go `*MmsValue`
  - `toGoValue()` returns `MmsDataAccessError` for DataAccessError type
- **Errors**
  - `GetMmsError()` – map C `MmsError` to Go errors for low-level MMS APIs
- **Server – connection and threadless**
  - `ClientConnection` – `PeerAddress()`, `LocalAddress()`, `SecurityToken()`, `Abort()`, `ClaimOwnership()`, `Release()`
  - `ConnectionIndicationHandler` and `SetConnectionIndicationHandler()` – connect/disconnect callbacks
  - `StartThreadless(port)`, `StopThreadless()`, `WaitReady(timeoutMs)`, `ProcessIncomingData()`, `PerformPeriodicTasks()`
- **Server – handlers and logging** (`server_handler.go`, `server_logging.go`)
  - Write access and control handlers, ACSE authenticator bridge, connection indication bridge
  - `LogStorageRef`, `NewLogStorageRef()`, `SetMaxLogEntries()`, logging API support
- **Server – MMS configuration and handlers** (`server_mms.go`, `shim.c`)
  - `SetFileAccessHandler()`, `InstallVariableListAccessHandler()`
  - `InstallReadJournalHandler()`, `InstallGetNameListHandler()`, `InstallObtainFileHandler()`, `InstallGetFileCompleteHandler()` (via C shim)
  - `SetMaxMmsConnections()`, `SetMaxConnections()`, `SetMaxMmsPduSize()`, `GetMaxMmsPduSize()`
  - `EnableMmsFileService()`, `EnableDynamicNamedVariableLists()`
  - `SetMaxAssociationSpecificDataSets()`, `SetMaxDomainSpecificDataSets()`, `SetMaxDataSetEntries()`, `EnableJournalService()`
  - `SetFilestoreBasepath()`, `GetFilestoreBasepath()`, `MmsServerConnection`
- **Documentation**
  - `GAPS.md` – MMS function coverage analysis (~98% coverage, client 100%)
  - `FUNCTIONS.md`, `STRUCTS.md`, `ENUMS.md` – API reference docs

### Changed

- GitHub Actions workflow (`build-libraries.yml`) updated for extra-functions branch
- `file_callback.c` replaced by `shim.c` (IedConnection getFile + MmsServer install* handler shims)

### Fixed

- **CGo memory management** – Fixed compilation warnings with Go 1.17+ strict type checking
  - Added `allocCString()` helper – allocator pattern for `C.CString()` with automatic cleanup
  - Added `allocCMalloc()` helper – allocator pattern for `C.malloc()` with automatic cleanup
  - Added `allocGo2CStr()` helper – allocator pattern for GB18030-encoded C strings with automatic cleanup
  - Replaced 174+ instances of `C.free(unsafe.Pointer(...))` across all Go files
  - Pattern: `cStr, free := allocCString("hello"); defer free()` eliminates unsafe pointer warnings
- **Type visibility** – Added CGo imports to `types.go` to resolve `MmsVariableSpecificationRef` undefined errors in language server

---

## [1.1.0] - 2026-02-01

### Breaking Changes

- **REQUIRED**: libiec61850 v1.6.1 (no backward compatibility with v1.5.3)
- **REQUIRED**: mbedtls dependency for R-GOOSE and R-SMV features
- Precompiled libraries rebuilt for all platforms with advanced features enabled
- Default server edition changed from Edition1 to Edition2

### Added

- **R-GOOSE** (Routable GOOSE over UDP/IP) support - enabled by default
- **R-SMV** (Routable Sampled Values over UDP/IP) support - enabled by default
- **SNTP Client** (Network time synchronization) - enabled by default
- **ClientGooseControlBlock wrapper** - Complete client-side GOOSE control block management
  - `GetGoCBValues()` - Read GOOSE control block configuration from server
  - `SetGoCBValues()` - Write GOOSE control block parameters with selective updates
  - `PhyComAddress` type for Ethernet/VLAN addressing (MAC, VLAN priority/ID, AppID)
  - GoCB element mask constants for selective parameter updates
- **macOS Intel (x86_64)** platform support with dedicated build configuration
- IEC 61850 Edition 2.1 support via `Edition21` constant in `server_config.go`
- `GetLibraryVersion()` function for runtime version checking
- `version_test.go` for version verification
- Automated rebuild script (`scripts/rebuild_libraries.sh`) with feature auto-configuration
- Comprehensive documentation in `COMPATIBILITY.md` - Feature coverage matrix
- Edition constants (`Edition1`, `Edition2`, `Edition21`) for type-safe configuration
- GitHub Actions workflow for automated multi-platform library builds

### Changed

- C headers synchronized to libiec61850 v1.6.1 from official release
- MMS header directory structure updated: `inc/mms/inc/` → `inc/mms/`
- All platform-specific config files (`config_*.go`) updated for new header paths
- Default edition in `NewServerConfig()` changed to `Edition2`
- Updated all documentation to reflect v1.6.1 requirements
- Rebuild script now auto-enables R-GOOSE, R-SMV, and SNTP features
- Enhanced error messages in version test with upgrade instructions

### Fixed

- **Windows x64 build issues** - Multiple fixes for MinGW linking and CGo configuration
  - Removed `-D__USE_MINGW_ANSI_STDIO=1` flag causing import stub conflicts
  - Added mbedtls platform-specific flags (`-DMBEDTLS_PLATFORM_SNPRINTF_ALT`, `-DMBEDTLS_PLATFORM_VSNPRINTF_ALT`)
  - Forced CGO to use MinGW gcc toolchain with explicit environment variables
  - Fixed shell environment (PowerShell vs MSYS2) for Go builds
- **macOS Intel (x86_64) build configuration** - Removed invalid `-Wl,-all_load` linker flags
- **test_getfile utility** - Corrected host:port parsing using `net.SplitHostPort` instead of `fmt.Sscanf`
- CGo include paths for all platforms (Linux x64/ARMv7l/ARMv8, macOS ARM64/Intel, Windows x64)
- Go package imports for header directories after rsync synchronization
- Duplicate package declarations in version-related files
- MMS header path references in all build configurations

### Documentation

- README.md: Added version requirements and build instructions
- All markdown files updated to reflect v1.6.1-only support
- Removed v1.5.3 backward compatibility references
- Added quick start guide for new users

### Migration Guide

Key migration steps from v1.0.13:
1. Update `go.mod`: `require github.com/boeboe/iec61850 v1.1.0`
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
- Config file model loading  M1/M2)
- macOS Intel x64
- Windows x64 (MinGW)
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

[1.1.1]: https://github.com/boeboe/iec61850/releases/tag/v1.1.1
[1.1.0]: https://github.com/boeboe/iec61850/releases/tag/v1.1.0
[1.0.13]: https://github.com/boeboe/iec61850/releases/tag/v1.0.13
