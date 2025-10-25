# Authentication

The CLI supports multiple authentication methods with the following precedence order:

1. **Environment Variables** (highest priority)
2. **Config File**
3. **Command-line Flags** (lowest priority)

## Environment Variables

Set your credentials as environment variables:

```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"
```

Add these to your shell profile (`.bashrc`, `.zshrc`, etc.) for persistence.

## Config File

Create a config file at `~/.trello-cli/config.yaml`:

```yaml
api_key: your-api-key
token: your-token
default_format: markdown
max_tokens: 4000
```

Or use the config command:

```bash
trello-cli config set --api-key "your-api-key" --token "your-token"
```

## Command-line Flags

Override credentials for specific commands:

```bash
trello-cli --api-key "your-api-key" --token "your-token" board list
```

## Testing Authentication

Test your setup with a simple command:

```bash
trello-cli board list
```

If authentication is working, you'll see your boards. If not, you'll get an authentication error.

## Troubleshooting

- **401 Unauthorized**: Check your API key and token
- **403 Forbidden**: Verify token permissions
- **Invalid key**: Ensure API key is correct and not expired
