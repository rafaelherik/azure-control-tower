# Publishing to Homebrew

There are two ways to make your package available via Homebrew:

## Option 1: Create a Homebrew Tap (Recommended)

A tap is a custom Homebrew repository. This is the easiest and fastest way to distribute your package.

### Step 1: Create a Tap Repository

1. Create a new GitHub repository named `homebrew-azct` (or `homebrew-azure-command-tower`)
   - Repository name must start with `homebrew-`
   - Make it public
   - Don't initialize with README, .gitignore, or license

2. Clone the repository:
   ```bash
   git clone https://github.com/rafaelherik/homebrew-azct.git
   cd homebrew-azct
   ```

### Step 2: Add the Formula

1. Copy your formula to the tap:
   ```bash
   cp /path/to/azure-control-tower/Formula/azct.rb /path/to/homebrew-azct/azct.rb
   ```

2. Update the formula URLs to point to your main repository:
   ```ruby
   class Azct < Formula
     desc "Terminal-based UI for exploring and managing Azure resources"
     homepage "https://github.com/rafaelherik/azure-control-tower"
     url "https://github.com/rafaelherik/azure-control-tower/archive/v0.0.1.tar.gz"
     sha256 "YOUR_SHA256_HERE"
     license "MIT"
     head "https://github.com/rafaelherik/azure-control-tower.git", branch: "main"

     depends_on "go" => :build

     def install
       system "go", "build", "-ldflags", "-s -w -X main.version=#{version}", "-o", bin/"azct", "./cmd/azct"
     end

     test do
       system "#{bin}/azct", "--version"
     end
   end
   ```

3. Calculate the SHA256 for the tarball:
   ```bash
   # Download the tarball first (or use the one from GitHub release)
   curl -L https://github.com/rafaelherik/azure-control-tower/archive/v0.0.1.tar.gz -o azure-control-tower-0.0.1.tar.gz
   shasum -a 256 azure-control-tower-0.0.1.tar.gz
   ```

4. Update the SHA256 in the formula

5. Commit and push:
   ```bash
   git add azct.rb
   git commit -m "Add azct formula"
   git push origin main
   ```

### Step 3: Install from Tap

Users can now install with:
```bash
brew tap rafaelherik/azct
brew install azct
```

Or in one command:
```bash
brew install rafaelherik/azct/azct
```

### Step 4: Update README

Add installation instructions to your main repository README:
```markdown
### Homebrew (macOS)

```bash
brew tap rafaelherik/azct
brew install azct
```
```

## Option 2: Submit to homebrew-core

This makes your package available via `brew install azct` without needing a tap, but requires approval.

### Prerequisites

- Package must be stable (not a pre-release)
- Must have at least 30 stars on GitHub
- Must have been tagged for at least 30 days
- Must have a stable release (not just HEAD)
- Must follow Homebrew naming conventions

### Steps

1. **Fork homebrew-core**:
   ```bash
   git clone https://github.com/Homebrew/homebrew-core.git
   cd homebrew-core
   ```

2. **Create formula**:
   ```bash
   # Copy your formula
   cp /path/to/azure-control-tower/Formula/azct.rb Formula/azct.rb
   ```

3. **Test locally**:
   ```bash
   brew install --build-from-source Formula/azct.rb
   brew test azct
   brew audit --strict azct
   ```

4. **Create pull request**:
   ```bash
   git checkout -b add-azct
   git add Formula/azct.rb
   git commit -m "azct 0.0.1 (new formula)"
   git push origin add-azct
   ```

5. **Submit PR** to https://github.com/Homebrew/homebrew-core

### Requirements for homebrew-core

- Formula must pass `brew audit --strict`
- Must have tests
- Must have a stable release (not just HEAD)
- Must follow Homebrew conventions
- May require discussion if name conflicts exist

## Automated Updates

### For Tap (Option 1)

You can automate formula updates in your tap repository:

1. Create a GitHub Actions workflow in your tap repo:
   ```yaml
   name: Update Formula
   on:
     repository_dispatch:
       types: [update-formula]
   jobs:
     update:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - name: Update formula
           run: |
             # Download new version
             # Calculate SHA256
             # Update formula
             # Commit and push
   ```

2. Trigger from your main repo's release workflow using `repository_dispatch`

### For homebrew-core (Option 2)

Updates are handled by Homebrew's automated bot or manual PRs from maintainers.

## Recommended Approach for v0.0.1

For your first release (v0.0.1), **use Option 1 (Tap)** because:
- ✅ Immediate availability
- ✅ Full control over updates
- ✅ No approval process
- ✅ Can update anytime
- ✅ Easier to iterate

You can always submit to homebrew-core later when:
- You have more stars/users
- Package is more mature
- You want broader distribution

## Quick Start: Create Tap Now

### Option A: Use the Setup Script

```bash
# From your azure-control-tower repository root
export GITHUB_USER=rafaelherik
export VERSION=0.0.1
./scripts/setup-homebrew-tap.sh
```

The script will:
1. Check if tap repository exists
2. Clone it if needed
3. Copy and update the formula
4. Calculate SHA256
5. Commit and push

### Option B: Manual Setup

```bash
# 1. Create tap repository on GitHub (via web UI)
# Name it: homebrew-azct
# Make it public, don't initialize with README

# 2. Clone it
git clone https://github.com/rafaelherik/homebrew-azct.git
cd homebrew-azct

# 3. Copy formula
cp ../azure-control-tower/Formula/azct.rb azct.rb

# 4. Update URLs in azct.rb (replace rafaelherik)

# 5. Calculate SHA256
curl -L https://github.com/rafaelherik/azure-control-tower/archive/v0.0.1.tar.gz -o /tmp/azure-control-tower-0.0.1.tar.gz
shasum -a 256 /tmp/azure-control-tower-0.0.1.tar.gz
# Update the SHA256 in azct.rb

# 6. Commit and push
git add azct.rb
git commit -m "Add azct formula v0.0.1"
git push origin main

# 7. Test installation
brew tap rafaelherik/azct
brew install azct
```

## Automated Updates

After creating your tap, you can set up automated updates:

1. **Add workflow to tap repository**: Copy `.github/workflows/homebrew-tap.yml.example` to your tap repo as `.github/workflows/update-formula.yml`

2. **Add secret to main repository**: 
   - Go to Settings → Secrets → Actions
   - Add secret: `HOMEBREW_TAP_REPO` = `rafaelherik/homebrew-azct`

3. **On each release**, the main repo will automatically trigger the tap to update

Alternatively, you can manually trigger updates:
```bash
gh api repos/rafaelherik/homebrew-azct/dispatches \
  --method POST \
  -f event_type=update-formula \
  -f client_payload='{"version":"0.0.2"}'
```

## Troubleshooting

### Formula not found
- Make sure tap repository is public
- Verify formula file is named correctly
- Check repository name starts with `homebrew-`

### SHA256 mismatch
- Recalculate SHA256 for the tarball
- Make sure you're using the correct tarball URL

### Build failures
- Test locally: `brew install --build-from-source Formula/azct.rb`
- Check Go version compatibility
- Verify all dependencies are available

