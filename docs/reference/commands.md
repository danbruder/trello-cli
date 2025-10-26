# Commands Overview

Trello CLI provides comprehensive commands for managing all Trello resources. All commands support the same global flags for authentication, output formatting, and LLM optimization.

## Available Commands

### Core Resource Commands

- **[boards](/reference/boards)** - Manage Trello boards
- **[lists](/reference/lists)** - Manage lists within boards  
- **[cards](/reference/cards)** - Manage cards within lists
- **[labels](/reference/labels)** - Manage labels and assign them to cards
- **[checklists](/reference/checklists)** - Manage checklists on cards
- **[members](/reference/members)** - View member information and boards
- **[attachments](/reference/attachments)** - Manage file attachments on cards

### Utility Commands

- **[batch](/reference/batch)** - Execute multiple operations from files or stdin
- **[config](/reference/config)** - Manage CLI configuration and credentials

## Global Flags

All commands support these global flags:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--api-key` | | Trello API key (overrides env/config) | |
| `--token` | | Trello token (overrides env/config) | |
| `--format` | `-f` | Output format (json, markdown) | markdown |
| `--fields` | | Specific fields to include in output | |
| `--max-tokens` | | Maximum tokens in output (0 = unlimited) | 0 |
| `--verbose` | `-v` | Verbose output | false |
| `--quiet` | `-q` | Quiet mode (minimal output) | false |
| `--debug` | | Debug mode (show API calls) | false |

## Command Structure

Most commands follow this pattern:

```bash
trlo <resource> <action> [arguments] [flags]
```

Examples:

```bash
# List resources
trlo board list
trlo card list --list <list-id>

# Get specific resource
trlo board get <board-id>
trlo card get <card-id>

# Create resources
trlo board create "My Board"
trlo card create --list <list-id> "My Card"

# Update resources
trlo card move <card-id> --list <target-list-id>
trlo label add <card-id> <label-id>
```

## LLM Optimization Features

### Field Filtering

Limit output to specific fields for token efficiency:

```bash
trlo card list --list <list-id> --fields name,desc,due
```

### Token Limits

Control output size for LLM context windows:

```bash
trlo board list --max-tokens 2000
```

### Output Formats

Choose between human-readable and machine-readable formats:

```bash
# Markdown (default)
trlo board list

# JSON for programmatic use
trlo board list --format json
```

## Error Handling

The CLI provides detailed error messages and appropriate exit codes:

- `0`: Success
- `1`: General error  
- `2`: Authentication error
- `3`: API error

## Getting Help

Get help for any command:

```bash
trlo --help
trlo board --help
trlo card create --help
```
