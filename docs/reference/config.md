# Configuration

Manage Trello CLI configuration including setting credentials and defaults.

## Commands

### `show`
Show current configuration settings.

```bash
trello-cli config show [flags]
```

**Examples:**
```bash
# Show current configuration
trello-cli config show

# Show configuration in JSON format
trello-cli config show --format json
```

### `set`
Set configuration values.

```bash
trello-cli config set [flags]
```

**Flags:**
- `--api-key` - Set the Trello API key
- `--token` - Set the Trello token
- `--default-format` - Set the default output format
- `--max-tokens` - Set the default maximum tokens

**Examples:**
```bash
# Set API credentials
trello-cli config set --api-key "your-api-key" --token "your-token"

# Set default format
trello-cli config set --default-format json

# Set maximum tokens
trello-cli config set --max-tokens 4000

# Set multiple values
trello-cli config set --api-key "key" --token "token" --default-format json --max-tokens 3000
```

### `path`
Show the path to the configuration file.

```bash
trello-cli config path [flags]
```

**Examples:**
```bash
# Show config file path
trello-cli config path

# Show path quietly
trello-cli config path --quiet
```

## Configuration File

The configuration file is stored at `~/.trello-cli/config.yaml` and has the following format:

```yaml
api_key: your-trello-api-key
token: your-trello-token
default_format: markdown  # or json
max_tokens: 4000         # 0 = unlimited
```

## Configuration Precedence

Configuration values are applied in the following order of precedence:

1. **Environment variables** (highest priority)
2. **Configuration file**
3. **Command-line flags** (lowest priority)

### Environment Variables

- `TRELLO_API_KEY` - Your Trello API key
- `TRELLO_TOKEN` - Your Trello access token

### Command-line Flags

- `--api-key` - Override API key
- `--token` - Override token
- `--format, -f` - Override output format
- `--max-tokens` - Override maximum tokens

## Common Use Cases

### Initial Setup
```bash
# Set up authentication
trello-cli config set --api-key "your-api-key" --token "your-token"

# Verify configuration
trello-cli config show
```

### Environment-Specific Configuration
```bash
# Development environment
trello-cli config set --default-format json --max-tokens 2000

# Production environment
trello-cli config set --default-format markdown --max-tokens 4000
```

### Configuration Management
```bash
#!/bin/bash
# Backup and restore configuration
CONFIG_PATH=$(trello-cli config path --quiet)

# Backup current config
cp "$CONFIG_PATH" ~/.trello-cli/config.yaml.backup

# Restore from backup
cp ~/.trello-cli/config.yaml.backup "$CONFIG_PATH"
```

### LLM Integration
```bash
# Set optimal defaults for LLM use
trello-cli config set --default-format json --max-tokens 3000

# Verify settings
trello-cli config show --format json
```
