# Quick Start: Building Libraries with GitHub Actions

This is a **step-by-step visual guide** for using GitHub Actions to build all platform libraries automatically.

## Prerequisites

- GitHub repository with Actions enabled (free for public repos)
- Push access to the repository
- ~15-30 minutes for the build to complete

## Step-by-Step Instructions

### 1. Navigate to Actions Tab

1. Open your repository on GitHub: `https://github.com/boeboe/iec61850`
2. Click the **Actions** tab at the top

### 2. Select the Workflow

1. In the left sidebar, click **Build libiec61850 Libraries**
2. On the right side, click the **Run workflow** dropdown button

### 3. Configure the Workflow Run

You'll see a form with options:

```
Use workflow from: [Branch: master ▼]

☐ Commit libraries to repository
   Commit built libraries to repository
```

**Configuration:**
- **Branch**: Select `master` (or `main` if that's your default)
- **Commit libraries**: ✅ **Check this box** to auto-commit libraries

### 4. Start the Build

Click the green **Run workflow** button.

You'll see a new workflow run appear with status "⚫ In progress".

### 5. Monitor Progress

Click on the workflow run to see details. You'll see 6 jobs running:

```
✓ Build Linux x86_64          (~5-10 min)
✓ Build Linux ARMv7l          (~8-12 min)
✓ Build Linux ARMv8 (ARM64)   (~8-12 min)
✓ Build macOS ARM64           (~10-15 min)
✓ Build Windows x64           (~10-15 min)
⏳ Verify and Commit Libraries (waits for all above)
```

### 6. View Build Summary

Once complete, scroll down to see the **Summary** with:

```markdown
# Build Summary

## Built Libraries

| Platform | Size | Status |
|----------|------|--------|
| Linux x86_64 | 2.4M | ✅ |
| Linux ARMv7l | 2.2M | ✅ |
| Linux ARMv8 | 2.3M | ✅ |
| macOS ARM64 | 2.5M | ✅ |
| Windows x64 | 2.6M | ✅ |

## Configuration
- libiec61850: v1.6.1
- R-GOOSE: ✅ Enabled
- R-SMV: ✅ Enabled
- SNTP Client: ✅ Enabled
```

### 7. Verify the Commit

If you enabled "Commit libraries", check your repository:

```bash
# Pull the latest changes
git pull origin master

# Verify libraries exist
ls -lh libiec61850/lib/*/libiec61850.a

# Check the commit
git log -1
```

You should see a commit like:
```
commit abc123...
Author: GitHub Actions <actions@github.com>
Date:   Sat Feb 1 12:34:56 2026 +0000

    Build: Update libiec61850 v1.6.1 libraries for all platforms
    
    Built with:
    - R-GOOSE enabled
    - R-SMV enabled
    - SNTP client enabled
    
    Platforms:
    - Linux x86_64
    - Linux ARMv7l
    - Linux ARMv8 (ARM64)
    - macOS ARM64
    - Windows x64
    
    Auto-built by GitHub Actions
```

### 8. Test Locally

```bash
# Test version
go test -v -run TestLibraryVersion

# Expected output:
# Using libiec61850 version: 1.6.1 (with R-GOOSE, R-SMV, SNTP enabled)

# Build project
go build ./...
```

### 9. Proceed with Release

Now you can create the v2.0.0 release:

```bash
# Tag the release
git tag -a v2.0.0 -m "Release v2.0.0"

# Push the tag
git push origin v2.0.0

# Create GitHub Release
# Go to: https://github.com/boeboe/iec61850/releases/new
# Select tag: v2.0.0
# Add release notes (see RELEASE_PROCESS.md)
```

## Alternative: Download Artifacts Manually

If you **didn't** enable "Commit libraries":

### 1. Download Artifacts

From the completed workflow run:
1. Scroll to **Artifacts** section at bottom
2. Download each platform:
   - `libiec61850-linux64.zip`
   - `libiec61850-linux-armv7l.zip`
   - `libiec61850-linux-armv8.zip`
   - `libiec61850-darwin-arm64.zip`
   - `libiec61850-windows-x64.zip`

### 2. Extract and Copy

```bash
# Extract each artifact
unzip libiec61850-linux64.zip
unzip libiec61850-linux-armv7l.zip
unzip libiec61850-linux-armv8.zip
unzip libiec61850-darwin-arm64.zip
unzip libiec61850-windows-x64.zip

# Copy to correct locations
cp libiec61850.a libiec61850/lib/linux64/
# ... repeat for other platforms
```

### 3. Commit Manually

```bash
git add libiec61850/lib/*/libiec61850.a
git commit -m "Update libiec61850 v1.6.1 libraries for all platforms"
git push origin master
```

## Troubleshooting

### Workflow Fails

**Check the logs:**
1. Click on the failed job (red ✗)
2. Expand each step to see errors
3. Common issues:
   - **mbedtls not found**: Package installation failed
   - **Cross-compiler missing**: apt-get failed for ARM tools
   - **CMake error**: Configuration issue

**Fix and re-run:**
1. Fix the workflow file if needed
2. Commit and push changes
3. Go to Actions → Re-run failed jobs

### Permission Denied on Commit

The workflow needs write access:
1. Go to **Settings** → **Actions** → **General**
2. Scroll to **Workflow permissions**
3. Select **Read and write permissions**
4. Click **Save**

### macOS Build Fails

Ensure you're using the M1 runner:
```yaml
runs-on: macos-14  # M1/M2 runner
```

Not:
```yaml
runs-on: macos-latest  # May be Intel
```

### Windows Build Fails with vcpkg Error

The workflow uses GitHub's pre-installed vcpkg. If it fails:
1. Check the error message
2. May need to specify vcpkg path explicitly
3. Or use Chocolatey to install mbedtls instead

## Best Practices

### For Production Releases

1. ✅ Always enable "Commit libraries"
2. ✅ Run from a clean branch (no uncommitted changes)
3. ✅ Verify artifacts before tagging release
4. ✅ Test locally after libraries are committed

### For Testing

1. ❌ Don't commit libraries (just download artifacts)
2. ✅ Test on a feature branch
3. ✅ Verify one platform at a time
4. ✅ Review logs for warnings

## Time Estimates

Typical workflow run times:

| Job | Time | Parallel |
|-----|------|----------|
| Build Linux x86_64 | 5-8 min | ✅ |
| Build Linux ARMv7l | 8-12 min | ✅ |
| Build Linux ARMv8 | 8-12 min | ✅ |
| Build macOS ARM64 | 10-15 min | ✅ |
| Build Windows x64 | 10-15 min | ✅ |
| Verify and Commit | 2-3 min | ⏳ Waits |
| **Total** | **15-30 min** | |

Jobs run in parallel, so total time is limited by the slowest job.

## Cost

- **Public repositories**: FREE (unlimited Actions minutes)
- **Private repositories**: Consumes GitHub Actions minutes
  - Free tier: 2,000 minutes/month
  - This workflow: ~80-120 minutes total (all platforms)

## Next Steps

After successful build:
1. ✅ Verify `go test -v -run TestLibraryVersion` passes
2. ✅ Update CHANGELOG.md with release notes
3. ✅ Tag release: `git tag -a v2.0.0`
4. ✅ Create GitHub Release
5. ✅ Update external tools to use v2.0.0

See [RELEASE_PROCESS.md](RELEASE_PROCESS.md) for complete release workflow.
