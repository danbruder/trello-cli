# Global Flags

All Trello CLI commands support these global flags for authentication, output formatting, and LLM optimization.

## Authentication Flags

### `--api-key`
Override the Trello API key for this command.

```bash
trello-cli board list --api-key "your-api-key"
```

### `--token`
Override the Trello token for this command.

```bash
trello-cli board list --token "your-token"
```

## Output Formatting

### `--format, -f`
Set the output format for the command.

**Values:**
- `markdown` (default) - Human-readable formatted output
- `json` - Structured JSON output

```bash
# Markdown output (default)
trello-cli board list

# JSON output
trello-cli board list --format json
trello-cli board list -f json
```

### `--fields`
Specify which fields to include in the output. Useful for reducing token usage.

```bash
# Only include name and description
trello-cli card list --list <list-id> --fields name,desc

# Include multiple fields
trello-cli board get <board-id> --fields name,desc,url,closed
```

### `--max-tokens`
Limit the output to a specific number of tokens. Set to `0` for unlimited.

```bash
# Limit to 2000 tokens
trello-cli board list --max-tokens 2000

# Unlimited tokens (default)
trello-cli board list --max-tokens 0
```

## Output Control

### `--verbose, -v`
Enable verbose output with additional details.

```bash
trello-cli board get <board-id> --verbose
trello-cli board get <board-id> -v
```

### `--quiet, -q`
Enable quiet mode with minimal output. Useful for scripting.

```bash
trello-cli card create --list <list-id> "New Card" --quiet
trello-cli card create --list <list-id> "New Card" -q
```

### `--debug`
Enable debug mode to show API calls and detailed information.

```bash
trello-cli board list --debug
```

## Flag Precedence

When multiple authentication methods are available, the precedence is:

1. Environment variables (`TRELLO_API_KEY`, `TRELLO_TOKEN`) - Highest priority
2. Configuration file (`~/.trello-cli/config.yaml`)
3. Command-line flags (`--api-key`, `--token`) - Lowest priority

## Examples

### Basic Usage
```bash
trello-cli board list
```

### JSON Output with Field Filtering
```bash
trello-cli card list --list <list-id> --format json --fields name,desc,due
```

### Quiet Mode for Scripting
```bash
trello-cli card create --list <list-id> "Task" --quiet
```

### Debug Mode for Troubleshooting
```bash
trello-cli board list --debug --verbose
```
