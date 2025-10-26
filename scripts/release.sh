#!/bin/bash

# Release script for trlo
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh 1.0.0

set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.0"
    exit 1
fi

# Validate version format (semantic versioning)
if [[ ! $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format X.Y.Z (e.g., 1.0.0)"
    exit 1
fi

echo "Releasing version $VERSION..."

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Error: Must be on main branch to release"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: Working directory is not clean"
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^v$VERSION$"; then
    echo "Error: Tag v$VERSION already exists"
    exit 1
fi

# Update version in main.go
sed -i.bak "s/version   = \".*\"/version   = \"$VERSION\"/" main.go
rm main.go.bak

# Update CHANGELOG.md
TODAY=$(date +'%Y-%m-%d')
sed -i.bak "s/## \[Unreleased\]/## \[Unreleased\]\n\n## \[$VERSION\] - $TODAY/" CHANGELOG.md
rm CHANGELOG.md.bak

# Commit changes
git add main.go CHANGELOG.md
git commit -m "Release version $VERSION"

# Create and push tag
git tag -a "v$VERSION" -m "Release version $VERSION"
git push origin main
git push origin "v$VERSION"

echo "Release $VERSION created successfully!"
echo ""
echo "Next steps:"
echo "1. Wait for GitHub Actions to build and create release"
echo "2. Update Homebrew formula with new version and SHA"
echo "3. Update Chocolatey package with new version and checksums"
echo "4. Test installation on all platforms"
echo ""
echo "GitHub release: https://github.com/danbruder/trello-cli/releases/tag/v$VERSION"
