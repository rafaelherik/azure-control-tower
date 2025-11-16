# Release Checklist for v0.0.1

## Pre-Release

- [ ] Update version in `VERSION` file
- [ ] Update version in `mkdocs.yml`
- [ ] Update version in `README.md`
- [ ] Update `CHANGELOG.md` with release notes
- [ ] Run all tests: `go test ./...`
- [ ] Build and test binaries for all platforms
- [ ] Update documentation if needed
- [ ] Set up Homebrew tap (see [Homebrew Guide](homebrew.md))
  - [ ] Create `homebrew-azct` repository on GitHub
  - [ ] Add formula to tap repository
  - [ ] Test installation: `brew tap yourusername/azct && brew install azct`

## Release Steps

1. **Create Git Tag**:
   ```bash
   git tag -a v0.0.1 -m "Release v0.0.1"
   git push origin v0.0.1
   ```

2. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create GitHub release
   - Upload binaries and checksums
   - Update Homebrew formula (if configured)

3. **Verify Release**:
   - Check GitHub releases page
   - Verify all binaries are uploaded
   - Update Homebrew tap (if not automated):
     ```bash
     # Use the setup script or manually update tap repository
     export GITHUB_USER=yourusername
     export VERSION=0.0.1
     ./scripts/setup-homebrew-tap.sh
     ```
   - Test Homebrew installation: `brew tap yourusername/azct && brew install azct`

## Post-Release

- [ ] Announce release (if applicable)
- [ ] Update documentation site (if not automatic)
- [ ] Monitor for issues
- [ ] Plan next release

## Manual Homebrew Formula Update (if needed)

If the Homebrew formula needs manual updates:

1. Calculate SHA256:
   ```bash
   shasum -a 256 azct-darwin-amd64
   ```

2. Update `Formula/azct.rb`:
   - Update version
   - Update SHA256
   - Update URL if needed

3. Test locally:
   ```bash
   brew install --build-from-source Formula/azct.rb
   ```

