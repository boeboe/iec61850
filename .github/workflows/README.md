# GitHub Actions Workflows

This directory contains automated workflows for building and managing the IEC 61850 Go library.

## Available Workflows

### build-libraries.yml

Automatically builds libiec61850 v1.6.1 static libraries for all supported platforms.

**Platforms Built:**
- Linux x86_64 (ubuntu-latest)
- Linux ARMv7l (cross-compiled on ubuntu-latest)
- Linux ARMv8/ARM64 (cross-compiled on ubuntu-latest)
- macOS ARM64 (macos-14 with M1 runner)
- Windows x64 (windows-latest with MSVC)

**Features Enabled:**
- R-GOOSE (Routable GOOSE over UDP/IP)
- R-SMV (Routable Sampled Values over UDP/IP)
- SNTP Client (Network time synchronization)

All builds include mbedtls for R-GOOSE and R-SMV support.

## Running the Build Workflow

### Manual Trigger (Recommended for Releases)

1. Go to **Actions** tab on GitHub
2. Select **Build libiec61850 Libraries**
3. Click **Run workflow**
4. Choose options:
   - **Branch**: Select your branch (usually `master` or `main`)
   - **Commit libraries**: Check this to automatically commit built libraries
5. Click **Run workflow**

The workflow will:
1. Build libraries for all 5 platforms in parallel
2. Upload each library as an artifact (downloadable for 30 days)
3. If "Commit libraries" is checked, commit them to the repository
4. Show a summary of all built libraries

### Automatic Trigger

The workflow also runs automatically when:
- Changes are pushed to the workflow file itself
- Changes are pushed to `scripts/rebuild_libraries.sh`

Note: Auto-triggered runs will **not** commit libraries by default.

## Workflow Details

### Build Process

Each platform job:
1. Checks out the repository
2. Installs platform-specific dependencies (cmake, mbedtls, cross-compilers)
3. Clones libiec61850 v1.6.1
4. Modifies `config/stack_config.h` to enable R-GOOSE, R-SMV, SNTP
5. Builds the static library
6. Uploads the library as an artifact

### Cross-Compilation

**Linux ARM builds** use cross-compilation:
- ARMv7l: `arm-linux-gnueabihf-gcc`
- ARMv8: `aarch64-linux-gnu-gcc`

mbedtls is also cross-compiled for ARM platforms.

### Artifacts

Built libraries are available as downloadable artifacts for 30 days:
- `libiec61850-linux64` - Linux x86_64
- `libiec61850-linux-armv7l` - Linux ARM 32-bit
- `libiec61850-linux-armv8` - Linux ARM 64-bit
- `libiec61850-darwin-arm64` - macOS Apple Silicon
- `libiec61850-windows-x64` - Windows x64

### Verification

After all platforms complete, a final job:
1. Downloads all artifacts
2. Copies them to correct locations
3. Runs `go build ./...` to verify Go compilation
4. Optionally commits libraries (if enabled)
5. Creates a build summary

## Using Built Libraries

### Download Artifacts Manually

1. Go to the workflow run in **Actions** tab
2. Scroll to **Artifacts** section
3. Download individual platform libraries
4. Extract and copy to `libiec61850/lib/<platform>/`

### Automatic Commit

If you enabled "Commit libraries" when running the workflow:
1. Libraries are automatically committed to your branch
2. Commit message includes build details
3. You can then merge/tag as needed for releases

## For Release v2.0.0

To rebuild all libraries for v2.0.0 release:

1. **Trigger the workflow**:
   ```
   Actions → Build libiec61850 Libraries → Run workflow
   ✓ Commit libraries to repository
   ```

2. **Wait for completion** (~15-30 minutes for all platforms)

3. **Verify the commit**:
   ```bash
   git pull
   ls -lh libiec61850/lib/*/libiec61850.a
   go test -v -run TestLibraryVersion
   ```

4. **Proceed with release**:
   ```bash
   git tag -a v2.0.0 -m "Release v2.0.0"
   git push origin v2.0.0
   ```

See [RELEASE_PROCESS.md](../../RELEASE_PROCESS.md) for complete release instructions.

## Troubleshooting

### Build Failures

**Check the workflow logs:**
- Click on the failed job
- Expand each step to see error details
- Common issues:
  - mbedtls installation failures (check package availability)
  - Cross-compiler not found (check apt install commands)
  - CMake configuration errors (check stack_config.h modifications)

**Windows-specific:**
- vcpkg integration issues: Ensure vcpkg path is correct
- MSVC not found: GitHub runner should have VS 2022 by default

**macOS-specific:**
- Use `macos-14` runner for M1/M2 support (not `macos-latest` which may be Intel)
- `sed -i ''` syntax for macOS (different from Linux)

### Artifact Downloads

Artifacts expire after 30 days. For long-term storage:
- Enable "Commit libraries" option
- Or download and commit manually

### Merge Conflicts

If auto-commit fails due to conflicts:
1. Download artifacts manually
2. Resolve conflicts locally
3. Commit and push

## Advanced Usage

### Modify Build Configuration

To change enabled features, edit the workflow file:

```yaml
env:
  LIBIEC61850_VERSION: v1.6.1
  ENABLE_R_GOOSE: 1        # Change to 0 to disable
  ENABLE_R_SMV: 1          # Change to 0 to disable
  ENABLE_SNTP: 1           # Change to 0 to disable
```

### Build Specific Platforms Only

Comment out jobs in the workflow:

```yaml
verify-and-commit:
  needs: [build-linux-amd64, build-macos-arm64]  # Only these platforms
```

### Custom CMake Options

Add to the cmake command in each job:

```yaml
cmake -DCMAKE_BUILD_TYPE=Release \
      -DYOUR_CUSTOM_OPTION=ON \
      ...
```

## Security

**Note**: The workflow uses:
- Public GitHub runners (no cost, but less control)
- Pre-installed tools (cmake, compilers, vcpkg)
- Official GitHub Actions (checkout@v4, upload-artifact@v4)

For production releases, verify:
- Artifact checksums match expected values
- Libraries are from trusted sources (official libiec61850 repo)
- No unexpected modifications in build process

## Cost

This workflow uses GitHub-hosted runners:
- **Free for public repositories** (unlimited minutes)
- Private repositories: Consumes GitHub Actions minutes

Typical run time:
- Linux jobs: ~5-10 minutes each
- macOS job: ~10-15 minutes
- Windows job: ~10-15 minutes
- Total: ~15-30 minutes (parallel execution)

## Support

For workflow issues:
- Check [GitHub Actions documentation](https://docs.github.com/en/actions)
- Review logs in the Actions tab
- Open an issue in the repository
