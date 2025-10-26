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

Create a config file at `~/.trlo/config.yaml`:

```yaml
api_key: your-api-key
token: your-token
default_format: markdown
max_tokens: 4000
```

Or use the config command:

```bash
trlo config set --api-key "your-api-key" --token "your-token"
```

## Command-line Flags

Override credentials for specific commands:

```bash
trlo --api-key "your-api-key" --token "your-token" board list
```

## Getting Your API Credentials

To use Trello CLI, you need to obtain API credentials from Trello:

1. Visit the [Trello Developer Portal](https://trello.com/app-key)
2. Copy your **API Key** (shown at the top of the page)
3. Click the "Token" link to generate a **Token** with appropriate permissions
4. Authorize the token (select read/write permissions as needed)
5. Copy the generated token

::: warning
Keep your API key and token secure. Never commit them to version control or share them publicly.
:::

## Testing Authentication

Test your setup with a simple command:

```bash
trlo board list
```

If authentication is working, you'll see your boards. If not, you'll get an authentication error.

## Troubleshooting

- **401 Unauthorized**: Check your API key and token
- **403 Forbidden**: Verify token permissions
- **Invalid key**: Ensure API key is correct and not expired
