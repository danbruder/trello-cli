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
trello-cli --help
```

You should see the help output with all available commands.

::: tip
If you built from source and haven't moved the binary to your PATH, you may need to use `./trello-cli --help` instead.
:::

## Next Steps

- [Set up authentication](/guide/authentication)
- [Try your first commands](/guide/quick-start)
