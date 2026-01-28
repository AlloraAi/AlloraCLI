# How to Create v2.0.0 Release

This document provides step-by-step instructions for creating the v2.0.0 release after the preparation PR is merged.

## Prerequisites

- All changes from the preparation PR have been merged to `main` branch
- You have write access to the repository
- You have git configured locally

## Steps to Create the Release

### 1. Switch to Main Branch and Pull Latest

```bash
# Switch to main branch
git checkout main

# Pull latest changes
git pull origin main

# Verify you have the latest changes
git log --oneline -5
```

You should see the commit "Prepare v2.0.0 release: Update CHANGELOG and add release notes"

### 2. Create the Git Tag

```bash
# Create an annotated tag for v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0 - Enhanced Platform Support and Installation Improvements"

# Verify the tag was created
git tag -l | grep v2.0.0
```

### 3. Push the Tag to GitHub

```bash
# Push the tag to GitHub
git push origin v2.0.0
```

This will trigger the GitHub Actions CI/CD pipeline which will:
- Run all tests
- Run security scans
- Build binaries for all platforms (Linux AMD64/ARM64, macOS Intel/ARM64, Windows AMD64)
- Create Docker images
- Create a GitHub Release
- Upload release assets

### 4. Monitor the Release Workflow

1. Go to: https://github.com/MohitSutharOfficial/AlloraCLI/actions
2. Look for the workflow run triggered by the tag
3. Monitor the progress of:
   - Test job
   - Lint job
   - Security scan job
   - Build job (for all platforms)
   - Docker build job
   - Release job

### 5. Verify the Release

Once the workflow completes successfully:

1. Go to: https://github.com/MohitSutharOfficial/AlloraCLI/releases
2. Verify the v2.0.0 release is created
3. Check that all release assets are uploaded:
   - `allora-linux-amd64.tar.gz`
   - `allora-linux-arm64.tar.gz`
   - `allora-darwin-amd64.tar.gz`
   - `allora-darwin-arm64.tar.gz`
   - `allora-windows-amd64.tar.gz`
4. Verify the release notes are included

### 6. Test the Release

Download and test the release on different platforms:

#### Linux/macOS
```bash
# Download the binary
curl -L https://github.com/MohitSutharOfficial/AlloraCLI/releases/download/v2.0.0/allora-linux-amd64 -o allora
chmod +x allora
./allora --version
# Should output: allora version v2.0.0 (commit: ..., date: ...)
```

#### Windows (PowerShell)
```powershell
# Test the automated installer
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/MohitSutharOfficial/AlloraCLI/main/scripts/install.ps1" -OutFile "$env:TEMP\install-allora.ps1"
& "$env:TEMP\install-allora.ps1"

# Verify
allora --version
```

### 7. Announce the Release (Optional)

Consider announcing the release on:
- GitHub Discussions
- Project README (if you have a "Latest Release" section)
- Social media
- Developer communities

## What Happens Automatically

When you push the v2.0.0 tag, GitHub Actions will automatically:

1. **Run Tests**: All Go tests and verification tests
2. **Security Scan**: CodeQL and vulnerability scanning
3. **Build Binaries**: For all supported platforms
4. **Create Release**: GitHub release with generated notes
5. **Upload Assets**: All binary archives
6. **Docker Images**: Build and push Docker images (if configured)

## Rollback Procedure (if needed)

If something goes wrong:

1. **Delete the tag locally**:
   ```bash
   git tag -d v2.0.0
   ```

2. **Delete the tag from GitHub**:
   ```bash
   git push origin :refs/tags/v2.0.0
   ```

3. **Delete the GitHub Release**:
   - Go to the releases page
   - Edit the v2.0.0 release
   - Delete the release

4. **Fix the issue** and repeat the process

## Troubleshooting

### Tag already exists
If you see "tag already exists":
```bash
# Delete and recreate
git tag -d v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0 --force
```

### Workflow fails
1. Check the Actions tab for error details
2. Fix the issue in the code
3. Delete the tag and release
4. Create a new commit with the fix
5. Create the tag again

### Release assets missing
1. Check the build job logs in GitHub Actions
2. Ensure all platform builds succeeded
3. Verify the release job uploaded the assets

## Post-Release Tasks

After successful release:

1. **Update Documentation**: Ensure all docs reference the correct version
2. **Update Package Managers** (future): Submit to Homebrew, Chocolatey, etc.
3. **Monitor Issues**: Watch for any issues reported with the new release
4. **Plan Next Release**: Start planning v2.1.0 features

## Release Checklist

Use this checklist when creating the release:

- [ ] Merged preparation PR to main
- [ ] Pulled latest changes from main
- [ ] Created git tag `v2.0.0`
- [ ] Pushed tag to GitHub
- [ ] Monitored CI/CD workflow
- [ ] Verified release created on GitHub
- [ ] Verified all release assets uploaded
- [ ] Tested installation on at least one platform
- [ ] Reviewed release notes
- [ ] Announced release (optional)

## Additional Resources

- [GitHub Releases Documentation](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

## Questions?

If you encounter any issues during the release process:
1. Check the GitHub Actions logs
2. Review this guide
3. Check the [CONTRIBUTING.md](CONTRIBUTING.md) guide
4. Open a discussion on GitHub

---

**Note**: This guide assumes you're creating the release from the main branch after merging the preparation PR. Always ensure the main branch is in a releasable state before creating a tag.
