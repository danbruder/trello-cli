# Installation

## Prerequisites

- Go 1.23 or later
- Trello API credentials (API key and token)

## Build from Source

```bash
git clone https://github.com/danbruder/trello-cli.git
cd trello-cli
go build -o trello-cli
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
