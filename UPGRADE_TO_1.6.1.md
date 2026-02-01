# Quick Upgrade Guide: libiec61850 v1.6.1

## TL;DR

This project now **requires libiec61850 v1.6.1** with R-GOOSE, R-SMV, and SNTP enabled.

### Quick Steps

1. **Install mbedtls** (required for R-GOOSE and R-SMV):
   ```bash
   # macOS
   brew install mbedtls
   
   # Ubuntu/Debian
   sudo apt-get install libmbedtls-dev
   
   # RHEL/CentOS
   sudo yum install mbedtls-devel
   ```

2. **Rebuild libraries**:
   ```bash
   cd /Users/bartvanbos/Git/boeboe/iec61850
   ./scripts/rebuild_libraries.sh
   ```

3. **Verify**:
   ```bash
   go test -v -run TestLibraryVersion
   # Expected: 1.6.1 (with R-GOOSE, R-SMV, SNTP enabled)
   ```

## What Changed

### Before (v1.5.3)
- Basic MMS, GOOSE L2, SV L2 support
- No routable protocols
- No SNTP client

### After (v1.6.1 - Required)
- ✅ All v1.5.3 features
- ✅ **R-GOOSE** - Routable GOOSE over UDP/IP (requires mbedtls)
- ✅ **R-SMV** - Routable Sampled Values over UDP/IP (requires mbedtls)
- ✅ **SNTP Client** - Network time synchronization
- ✅ IEC 61850 Edition 2.1 support

## Build Configuration

The rebuild script automatically configures:
```c
#define CONFIG_IEC61850_R_GOOSE 1       // ENABLED
#define CONFIG_IEC61850_R_SMV 1         // ENABLED
#define CONFIG_IEC61850_SNTP_CLIENT 1   // ENABLED
```

## Platform Support

The rebuild script supports:
- ✅ **macOS ARM64** (M1/M2) - native build
- ✅ **Linux x86_64** - native build
- ✅ **Linux ARMv7l** - cross-compile (requires toolchain)
- ✅ **Linux ARMv8** - cross-compile (requires toolchain)
- ⚠️ **Windows x64** - manual build required (see REBUILD_LIBRARIES.md)

## Rebuild Script Details

The script will:
1. Clone libiec61850 v1.6.1 from official repository
2. Enable R-GOOSE, R-SMV, and SNTP in `config/stack_config.h`
3. Build static libraries for your platform
4. Copy libraries to `libiec61850/lib/<platform>/`
5. Clean up temporary files

Estimated time: 5-10 minutes per platform

## Testing

After rebuilding:
```bash
# Version check
go test -v -run TestLibraryVersion

# Full build
go build ./...

# Run test suite
go test -v ./test/...
```

## Troubleshooting

### "mbedtls not found"
```bash
# Verify mbedtls installation
pkg-config --modversion mbedtls

# macOS
brew install mbedtls

# Linux
sudo apt-get install libmbedtls-dev
```

### "Version mismatch: Expected 1.6.1, got 1.5.3"
You need to rebuild the libraries. The headers are v1.6.1 but precompiled libraries are still v1.5.3.
Run: `./scripts/rebuild_libraries.sh`

### Cross-compilation fails
Install cross-compiler toolchains:
```bash
# ARMv7l
sudo apt-get install gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf

# ARMv8
sudo apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

## Documentation

- [REBUILD_LIBRARIES.md](REBUILD_LIBRARIES.md) - Detailed build instructions
- [COMPATIBILITY.md](COMPATIBILITY.md) - Feature coverage matrix
- [README.md](README.md) - Project overview and usage

## Support

For issues with:
- **Library build**: See [REBUILD_LIBRARIES.md](REBUILD_LIBRARIES.md)
- **Feature compatibility**: See [COMPATIBILITY.md](COMPATIBILITY.md)
- **Go bindings**: Check test examples in `test/` directory
- **libiec61850 itself**: https://github.com/mz-automation/libiec61850

## Next Steps After Upgrade

Once v1.6.1 libraries are built:

1. **Add R-GOOSE Go bindings** (currently C-only):
   - Create `r_goose_publisher.go`
   - Create `r_goose_subscriber.go`
   - Wrap R-Session protocol APIs

2. **Add R-SMV Go bindings** (currently C-only):
   - Create `r_smv_publisher.go`
   - Create `r_smv_subscriber.go`

3. **Add SNTP Client Go bindings**:
   - Create `sntp_client.go`
   - Wrap SNTP client APIs

4. **Test with real devices**:
   - Verify R-GOOSE interoperability
   - Test R-SMV data exchange
   - Validate SNTP time sync

See [COMPATIBILITY.md](COMPATIBILITY.md) for the complete missing bindings list.
