# Installation

## Prerequisites

- Trello API credentials (API key and token)

## Package Managers

### Homebrew (macOS)

The easiest way to install on macOS:

```bash
brew tap danbruder/tap
brew install trello-cli
```

### Chocolatey (Windows)

For Windows users:

```powershell
choco install trello-cli
```

## Manual Installation

### Download Pre-built Binaries

1. Go to the [Releases page](https://github.com/danbruder/trello-cli/releases)
2. Download the appropriate binary for your platform:
   - `trello-cli-linux-amd64` for Linux x86_64
   - `trello-cli-linux-arm64` for Linux ARM64
   - `trello-cli-darwin-amd64` for macOS Intel
   - `trello-cli-darwin-arm64` for macOS Apple Silicon
   - `trello-cli-windows-amd64.exe` for Windows x86_64
   - `trello-cli-windows-arm64.exe` for Windows ARM64

3. Make it executable and move to your PATH:

```bash
# Linux/macOS
chmod +x trello-cli-*
sudo mv trello-cli-* /usr/local/bin/trello-cli

# Windows
# Move the .exe file to a directory in your PATH
```

### Docker

Run directly from Docker:

```bash
docker run --rm -it ghcr.io/danbruder/trello-cli:latest
```

### Build from Source

For development or if you need the latest features:

```bash
git clone https://github.com/danbruder/trello-cli.git
cd trello-cli
go build -o trello-cli .
```

## Get Trello API Credentials

1. Go to [Trello Developer Portal](https://trello.com/app-key)
2. Copy your API Key
3. Generate a token with appropriate permissions
4. Set up authentication (see [Authentication Guide](/guide/authentication))

## Verify Installation

```bash
./trello-cli --help
```

You should see the help output with all available commands.

## Next Steps

- [Set up authentication](/guide/authentication)
- [Try your first commands](/guide/quick-start)
