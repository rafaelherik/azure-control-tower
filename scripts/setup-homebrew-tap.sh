#!/bin/bash

# Script to set up Homebrew tap for azct
# Usage: ./scripts/setup-homebrew-tap.sh

set -e

REPO_NAME="azure-control-tower"
TAP_NAME="homebrew-azct"
GITHUB_USER="${GITHUB_USER:-yourusername}"
VERSION="${VERSION:-0.0.1}"

echo "Setting up Homebrew tap for azct"
echo "================================="
echo ""

# Check if tap repo exists
echo "Step 1: Checking if tap repository exists..."
if [ -d "../${TAP_NAME}" ]; then
    echo "✓ Tap repository already cloned"
    cd "../${TAP_NAME}"
else
    echo "✗ Tap repository not found locally"
    echo ""
    echo "Please create the tap repository on GitHub first:"
    echo "  1. Go to https://github.com/new"
    echo "  2. Repository name: ${TAP_NAME}"
    echo "  3. Make it public"
    echo "  4. Don't initialize with README"
    echo ""
    read -p "Press Enter after creating the repository..."
    
    echo "Cloning tap repository..."
    git clone "https://github.com/${GITHUB_USER}/${TAP_NAME}.git" "../${TAP_NAME}"
    cd "../${TAP_NAME}"
fi

# Copy formula
echo ""
echo "Step 2: Copying formula..."
if [ -f "../${REPO_NAME}/Formula/azct.rb" ]; then
    cp "../${REPO_NAME}/Formula/azct.rb" "azct.rb"
    echo "✓ Formula copied"
else
    echo "✗ Formula not found at ../${REPO_NAME}/Formula/azct.rb"
    exit 1
fi

# Update URLs in formula
echo ""
echo "Step 3: Updating URLs in formula..."
sed -i.bak "s|homepage \".*\"|homepage \"https://github.com/${GITHUB_USER}/${REPO_NAME}\"|" azct.rb
sed -i.bak "s|url \".*\"|url \"https://github.com/${GITHUB_USER}/${REPO_NAME}/archive/v${VERSION}.tar.gz\"|" azct.rb
sed -i.bak "s|head \".*\"|head \"https://github.com/${GITHUB_USER}/${REPO_NAME}.git\", branch: \"main\"|" azct.rb
rm azct.rb.bak
echo "✓ URLs updated"

# Calculate SHA256
echo ""
echo "Step 4: Calculating SHA256..."
echo "Downloading tarball..."
curl -L -s "https://github.com/${GITHUB_USER}/${REPO_NAME}/archive/v${VERSION}.tar.gz" -o "/tmp/azct-${VERSION}.tar.gz"
SHA256=$(shasum -a 256 "/tmp/azct-${VERSION}.tar.gz" | cut -d' ' -f1)
echo "SHA256: ${SHA256}"

# Update SHA256 in formula
sed -i.bak "s|sha256 \".*\"|sha256 \"${SHA256}\"|" azct.rb
rm azct.rb.bak
echo "✓ SHA256 updated"

# Show formula
echo ""
echo "Step 5: Formula preview:"
echo "========================"
cat azct.rb
echo ""

# Commit
echo ""
read -p "Commit and push? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git add azct.rb
    git commit -m "Add azct formula v${VERSION}"
    git push origin main
    echo ""
    echo "✓ Formula pushed to tap repository"
    echo ""
    echo "Users can now install with:"
    echo "  brew tap ${GITHUB_USER}/azct"
    echo "  brew install azct"
else
    echo "Formula ready but not committed. Review azct.rb and commit manually."
fi

