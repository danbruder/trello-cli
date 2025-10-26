# Distribution Guide

This document outlines the complete distribution strategy for the Trello CLI tool.

## Distribution Files

All distribution-related files are organized in the `dist/` directory:
- `dist/chocolatey/` - Chocolatey package files
- `dist/Formula/` - Homebrew formula
- `dist/debian/` - Debian package files
- `dist/rpm/` - RPM package files

## Distribution Channels

### 1. GitHub Releases (Primary)
- **Target**: All platforms
- **Method**: Automated via GitHub Actions
- **Files**: Cross-platform binaries, checksums, release notes
- **Update**: Automatic on version tags

### 2. Homebrew (macOS)
- **Target**: macOS users
- **Method**: Custom tap repository
- **Installation**: `brew install danbruder/tap/trello-cli`
- **Update**: Manual formula updates

### 3. Chocolatey (Windows)
- **Target**: Windows users
- **Method**: Package submission to Chocolatey
- **Installation**: `choco install trello-cli`
- **Update**: Manual package updates

### 4. APT (Debian/Ubuntu)
- **Target**: Debian-based Linux distributions
- **Method**: Custom APT repository
- **Installation**: Repository setup + `apt install trello-cli`
- **Update**: Repository updates

### 5. YUM/DNF (Red Hat/CentOS/Fedora)
- **Target**: Red Hat-based Linux distributions
- **Method**: Custom RPM repository
- **Installation**: Repository setup + `dnf install trello-cli`
- **Update**: Repository updates

### 6. Docker Hub/GitHub Container Registry
- **Target**: Containerized environments
- **Method**: Multi-architecture Docker images
- **Installation**: `docker pull ghcr.io/danbruder/trello-cli`
- **Update**: Automated on releases

## Release Process

### 1. Version Management
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Update version in multiple places:
  - `main.go` (version variable)
  - `CHANGELOG.md`
  - Package manifests
  - GitHub release notes

### 2. Automated Release Workflow
1. Create and push version tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions automatically:
   - Builds binaries for all platforms
   - Creates GitHub release
   - Builds and pushes Docker images
   - Generates checksums

### 3. Manual Package Updates
After automated release:
1. Update Homebrew formula with new version and SHA
2. Update Chocolatey package with new version and checksums
3. Update APT/RPM repositories (when implemented)

## Build Configuration

### Cross-Platform Builds
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, arm64

### Build Flags
- `-ldflags="-s -w"`: Strip debug info and symbol table
- `-X main.version=$(VERSION)`: Embed version information
- `CGO_ENABLED=0`: Static linking for better compatibility

## Package Repository Setup

### Homebrew Tap
1. ✅ Repository created: `github.com/danbruder/homebrew-tap`
2. ✅ Formula added: `dist/Formula/trello-cli.rb`
3. ✅ Installation: `brew tap danbruder/tap && brew install trello-cli`

### Chocolatey Package
1. Create package structure with `.nuspec` and install script
2. Submit to Chocolatey community repository
3. Maintain package updates

### APT Repository (Future)
1. Set up signing keys
2. Create repository structure
3. Implement automated package building
4. Set up repository hosting

### RPM Repository (Future)
1. Set up signing keys
2. Create repository structure
3. Implement automated package building
4. Set up repository hosting

## Security Considerations

### Code Signing
- **macOS**: Notarization for Gatekeeper compatibility
- **Windows**: Authenticode signing for SmartScreen
- **Linux**: GPG signing for package verification

### Binary Verification
- SHA256 checksums for all releases
- GPG signatures for packages
- Reproducible builds where possible

## Monitoring and Analytics

### Download Tracking
- GitHub release download counts
- Package manager installation metrics
- Docker Hub pull statistics

### User Feedback
- GitHub Issues for bug reports
- GitHub Discussions for feature requests
- Community forums for support

## Maintenance Tasks

### Regular Updates
- Monitor for security vulnerabilities
- Update dependencies
- Test on all supported platforms
- Update documentation

### Release Schedule
- **Major releases**: Breaking changes, new features
- **Minor releases**: New features, backward compatible
- **Patch releases**: Bug fixes, security updates

## Future Enhancements

### Additional Package Managers
- Snap (Ubuntu)
- Flatpak (Linux)
- Scoop (Windows)
- MacPorts (macOS)

### Distribution Improvements
- Automated package repository updates
- Code signing implementation
- Reproducible builds
- Automated testing on all platforms

### Community Contributions
- Package maintainer guidelines
- Contribution templates
- Automated issue labeling
- Community moderation tools
