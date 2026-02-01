# libiec61850 v1.6.1 Feature Coverage

This document describes the compatibility between the official libiec61850 C library (v1.6.1) and these Go bindings.

## Version Information

- **libiec61850 Version (Required)**: 1.6.1
- **Go Bindings Version**: libiec61850 v1.6.1 only
- **Build Configuration**: R-GOOSE, R-SMV, and SNTP client enabled by default
- **Last Updated**: February 2026

**Note**: This project requires libiec61850 v1.6.1 with advanced features enabled:
- ✅ R-GOOSE (Routable GOOSE over UDP/IP)
- ✅ R-SMV (Routable Sampled Values over UDP/IP)
- ✅ SNTP Client (Network time synchronization)

To build the required libraries, run `./scripts/rebuild_libraries.sh`.

To check the runtime version:
```go
version := iec61850.GetLibraryVersion()
fmt.Printf("Using libiec61850 version: %s\n", version)
```

## Core Features

| Feature | C Library | Go Bindings | Notes |
|---------|-----------|-------------|-------|
| **MMS Client** | ✅ | ✅ | Full support for read/write operations |
| **MMS Server** | ✅ | ✅ | Full support with callbacks |
| **Client Directory Services** | ✅ | ✅ | Modern API (GetServerDirectory, GetDataDirectoryWithFC) |
| **Report Control Blocks (RCB)** | ✅ | ✅ | Buffered and unbuffered reports |
| **Data Sets** | ✅ | ✅ | Dynamic and static data sets |
| **Setting Groups** | ✅ | ✅ | Setting group control via client/server |
| **Control Services** | ✅ | ✅ | Direct/SBO control (SPC, DPC, APC, INC) |
| **File Services** | ✅ | ✅ | GetFile/SetFile operations |
| **TLS Support** | ✅ | ✅ | Mutual TLS authentication for MMS |
| **IEC 61850 Edition 1** | ✅ | ✅ | Full support |
| **IEC 61850 Edition 2** | ✅ | ✅ | Full support (default) |
| **IEC 61850 Edition 2.1** | ✅ | ✅ | Use `iec61850.Edition21` constant |

## GOOSE (Generic Object Oriented Substation Event)

| Feature | C Library | Go Bindings | Notes |
|---------|-----------|-------------|-------|
| **GOOSE Publisher (L2)** | ✅ | ✅ | Ethernet Layer 2 GOOSE |
| **GOOSE Subscriber (L2)** | ✅ | ✅ | Platform-specific (Linux AMD64) |
| **GOOSE Control Block (Client)** | ✅ | ✅ | Read/Write GoCB via GetGoCBValues/SetGoCBValues |
| **R-GOOSE (Routable)** | ✅ | ✅ | Enabled by default (requires mbedtls) |

### R-GOOSE Support

R-GOOSE (Routable GOOSE over UDP/IP) is **enabled by default** in this project.

**Prerequisites**:
- mbedtls library must be installed:
  ```bash
  # macOS
  brew install mbedtls
  
  # Ubuntu/Debian
  sudo apt-get install libmbedtls-dev
  
  # RHEL/CentOS
  sudo yum install mbedtls-devel
  ```

**Note**: Go bindings for R-GOOSE API require manual wrapper implementation (see [Missing Go Bindings](#known-limitations)).

## Sampled Values (SV)

| Feature | C Library | Go Bindings | Notes |
|---------|-----------|-------------|-------|
| **SV Publisher (L2)** | ✅ | ✅ | IEC 61850-9-2 Layer 2 |
| **SV Subscriber (L2)** | ✅ | ✅ | Platform-specific (Linux AMD64) |
| **R-SMV (Routable)** | ✅ | ✅ | Enabled by default (requires mbedtls) |

### R-SMV Support

R-SMV (Routable Sampled Values over UDP/IP) is **enabled by default** in this project.

**Prerequisites**:
- mbedtls library must be installed (same as R-GOOSE)

**Note**: Go bindings for R-SMV API require manual wrapper implementation (see [Missing Go Bindings](#known-limitations)).

## Time Synchronization

| Feature | C Library | Go Bindings | Notes |
|---------|-----------|-------------|-------|
| **SNTP Client** | ✅ | ✅ | Enabled by default |
| **System Time Callbacks** | ✅ | ✅ | Custom time source via callbacks |

### SNTP Client Support

The SNTP (Simple Network Time Protocol) client is **enabled by default** in this project.

**Go bindings example** (requires implementation in new file `sntp_client.go`):
```go
package iec61850

/*
#include "sntp_client.h"
*/
import "C"

type SNTPClient struct {
    client C.SNTPClient
}

func NewSNTPClient() *SNTPClient {
    return &SNTPClient{
        client: C.SNTPClient_create(),
    }
}
```

## SCL (System Configuration Language)

| Feature | C Library | Go Bindings | Notes |
|---------|-----------|-------------|-------|
| **ICD File Parser** | ✅ | ✅ | Parse .icd/.cid files via scl package |
| **Static Model Generator** | ✅ | ✅ | Generate C model code from SCL |
| **Config File Parser** | ✅ | ✅ | Load .cfg model files |

## Platform Support

| Platform | Architecture | Status | Library Path |
|----------|-------------|--------|--------------|
| Linux | x86_64 | ✅ | libiec61850/lib/linux64/ |
| Linux | ARMv7l | ✅ | libiec61850/lib/linux_armv7l/ |
| Linux | ARMv8 (ARM64) | ✅ | libiec61850/lib/linux_armv8/ |
| macOS | ARM64 (M1/M2) | ✅ | libiec61850/lib/darwin_armv8/ |
| macOS | Intel x64 | ✅ | libiec61850/lib/darwin_amd64/ |
| Windows | x86_64 (MinGW) | ✅ | libiec61850/lib/win64/ |

## Known Limitations

### Missing Go Bindings (C library supports, but not wrapped)

1. **Async GOOSE Control**:
   - `IedConnection_getGoCBValuesAsync()` with callback handlers
   - `IedConnection_setGoCBValuesAsync()` with callback handlers

2. **Log Service**:
   - `IedConnection_queryLogAfter()`
   - `IedConnection_queryLogByTime()`

3. **Advanced MMS Services**:
   - Journal/Log reading APIs
   - Some MMS type introspection functions

### Enabled Features (v1.6.1 Build Configuration)

The following advanced features are **enabled by default** in the v1.6.1 precompiled libraries:

- ✅ R-GOOSE (Routable GOOSE over UDP/IP) - requires mbedtls
- ✅ R-SMV (Routable Sampled Values over UDP/IP) - requires mbedtls
- ✅ SNTP Client (Network time synchronization)

These are configured in `stack_config.h` during library rebuild.

## Rebuilding Precompiled Libraries

If you need to enable disabled features or update to a newer libiec61850 version:

Quick rebuild workflow:
```bash
# Clone official libiec61850
git clone --branch v1.6.1 https://github.com/mz-automation/libiec61850.git

# Configure stack features
cd libiec61850
# Edit config/stack_config.h to enable desired features

# Build for your platform
mkdir build && cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
make -j$(nproc)

# Copy library to Go bindings
cp output/libiec61850.a /path/to/iec61850/libiec61850/lib/<platform>/
```

## Testing Compatibility

Run the version test to verify your bindings are using v1.6.1:
```bash
cd /Users/bartvanbos/Git/boeboe/iec61850
go test -v -run TestLibraryVersion
```

Run full test suite:
```bash
go test -v ./test/...
```

## Contributing

When adding new Go wrappers for missing C API functions:

1. Follow the existing CGo patterns in client.go, server.go, etc.
2. Always use `defer C.free(unsafe.Pointer(cString))` after `C.CString()`
3. Add comprehensive tests in `test/<feature>/`
4. Update this compatibility matrix

## References

- [libiec61850 Official Repository](https://github.com/mz-automation/libiec61850)
- [IEC 61850 Standard Overview](https://en.wikipedia.org/wiki/IEC_61850)
- [Project README](README.md)
- [Copilot AI Instructions](.github/copilot-instructions.md)
