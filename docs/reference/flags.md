# Global Flags

All Trello CLI commands support these global flags for authentication, output formatting, and LLM optimization.

## Authentication Flags

### `--api-key`
Override the Trello API key for this command.

```bash
trlo board list --api-key "your-api-key"
```

### `--token`
Override the Trello token for this command.

```bash
trlo board list --token "your-token"
```

## Output Formatting

### `--format, -f`
Set the output format for the command.

**Values:**
- `markdown` (default) - Human-readable formatted output
- `json` - Structured JSON output

```bash
# Markdown output (default)
trlo board list

# JSON output
trlo board list --format json
trlo board list -f json
```

### `--fields`
Specify which fields to include in the output. Useful for reducing token usage.

```bash
# Only include name and description
trlo card list --list <list-id> --fields name,desc

# Include multiple fields
trlo board get <board-id> --fields name,desc,url,closed
```

### `--max-tokens`
Limit the output to a specific number of tokens. Set to `0` for unlimited.

```bash
# Limit to 2000 tokens
trlo board list --max-tokens 2000

# Unlimited tokens (default)
trlo board list --max-tokens 0
```

## Output Control

### `--verbose, -v`
Enable verbose output with additional details.

```bash
trlo board get <board-id> --verbose
trlo board get <board-id> -v
```

### `--quiet, -q`
Enable quiet mode with minimal output. Useful for scripting.

```bash
trlo card create --list <list-id> "New Card" --quiet
trlo card create --list <list-id> "New Card" -q
```

### `--debug`
Enable debug mode to show API calls and detailed information.

```bash
trlo board list --debug
```

## Flag Precedence

When multiple authentication methods are available, the precedence is:

1. Environment variables (`TRELLO_API_KEY`, `TRELLO_TOKEN`) - Highest priority
2. Configuration file (`~/.trlo/config.yaml`)
3. Command-line flags (`--api-key`, `--token`) - Lowest priority

## Examples

### Basic Usage
```bash
trlo board list
```

### JSON Output with Field Filtering
```bash
trlo card list --list <list-id> --format json --fields name,desc,due
```

### Quiet Mode for Scripting
```bash
trlo card create --list <list-id> "Task" --quiet
```

### Debug Mode for Troubleshooting
```bash
trlo board list --debug --verbose
```
