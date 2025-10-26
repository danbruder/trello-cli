# Installation

## Prerequisites

- Trello API credentials (API key and token)

## Package Managers

### Homebrew (macOS)

The easiest way to install on macOS:

```bash
brew tap danbruder/tap
brew install trlo
```

### Chocolatey (Windows)

For Windows users:

```powershell
choco install trlo
```

## Manual Installation

### Download Pre-built Binaries

1. Go to the [Releases page](https://github.com/danbruder/trlo/releases)
2. Download the appropriate binary for your platform:
   - `trlo-linux-amd64` for Linux x86_64
   - `trlo-linux-arm64` for Linux ARM64
   - `trlo-darwin-amd64` for macOS Intel
   - `trlo-darwin-arm64` for macOS Apple Silicon
   - `trlo-windows-amd64.exe` for Windows x86_64
   - `trlo-windows-arm64.exe` for Windows ARM64

3. Make it executable and move to your PATH:

```bash
# Linux/macOS
chmod +x trlo-*
sudo mv trlo-* /usr/local/bin/trlo

# Windows
# Move the .exe file to a directory in your PATH
```

### Docker

Run directly from Docker:

```bash
# With environment variables
docker run --rm -it \
  -e TRELLO_API_KEY="your-api-key" \
  -e TRELLO_TOKEN="your-token" \
  ghcr.io/danbruder/trlo:latest board list

# Or with a config file
docker run --rm -it \
  -v ~/.trlo:/root/.trlo \
  ghcr.io/danbruder/trlo:latest board list
```

### Build from Source

For development or if you need the latest features:

```bash
git clone https://github.com/danbruder/trlo.git
cd trlo
go build -o trlo .
```

## Get Trello API Credentials

1. Go to [Trello Developer Portal](https://trello.com/app-key)
2. Copy your API Key
3. Generate a token with appropriate permissions
4. Set up authentication (see [Authentication Guide](/guide/authentication))

## Verify Installation

```bash
trlo --help
```

You should see the help output with all available commands.

::: tip
If you built from source and haven't moved the binary to your PATH, you may need to use `./trlo --help` instead.
:::

## Next Steps

- [Set up authentication](/guide/authentication)
- [Try your first commands](/guide/quick-start)
