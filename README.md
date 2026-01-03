![logo](./logo-simple.png)

# Trello CLI for humans and LLMs

A comprehensive Trello CLI tool built in Go that provides full access to Trello's API with features optimized for LLM integration including context optimization, batch operations, and flexible output formats.

## Features

- **Comprehensive Operations**: Full CRUD operations on boards, lists, cards, labels, checklists, members, and attachments
- **LLM-Optimized Output**: Both Markdown and JSON formats with token counting and field filtering
- **Flexible Authentication**: Environment variables, config file, or command-line flags with precedence
- **Batch Operations**: Execute multiple operations from files or stdin for automation
- **Context Optimization**: Token limits, field filtering, and summarization for LLM use cases
- **Scripting Support**: Designed for automation and integration with LLM workflows

## Installation

### Prerequisites

- Trello API credentials (API key and token)

### Package Managers

#### Homebrew (macOS)

```bash
brew tap danbruder/tap
brew install trello-cli
```

### Build from Source

```bash
git clone https://github.com/danbruder/trello-cli.git
cd trello-cli
go build -o trello-cli .
```

### Get Trello API Credentials

Check out the authorization docs [here](https://trello-cli.netlify.app/guide/authentication.html).

## Authentication

The CLI supports multiple authentication methods with the following precedence order:

1. **Environment Variables** (highest priority)
2. **Config File**
3. **Command-line Flags** (lowest priority)

### Environment Variables

```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"
```

### Config File

Create a config file at `~/.trello-cli/config.yaml`:

```yaml
api_key: your-api-key
token: your-token
default_format: json
max_tokens: 4000
```

Or use the config command:

```bash
trello-cli config set --api-key "your-api-key" --token "your-token"
```

### Command-line Flags

```bash
trello-cli --api-key "your-api-key" --token "your-token" board list
```

## Usage

### Basic Commands

#### Boards

```bash
# List all boards
trello-cli board list

# Get board details
trello-cli board get <board-id>

# Create a new board
trello-cli board create "My New Board"

# Delete a board
trello-cli board delete <board-id>

# Add member to board
trello-cli board add-member <board-id> user@example.com
```

#### Lists

```bash
# List all lists on a board
trello-cli list list --board <board-id>

# Get list details
trello-cli list get <list-id>

# Create a new list
trello-cli list create --board <board-id> "New List"

# Archive a list
trello-cli list archive <list-id>
```

#### Cards

```bash
# List all cards in a list
trello-cli card list --list <list-id>

# Get card details
trello-cli card get <card-id>

# Create a new card
trello-cli card create --list <list-id> "New Card"

# Move a card to another list
trello-cli card move <card-id> --list <target-list-id>

# Copy a card
trello-cli card copy <card-id> --list <target-list-id>

# Archive a card
trello-cli card archive <card-id>

# Delete a card
trello-cli card delete <card-id>
```

#### Labels

```bash
# List all labels on a board
trello-cli label list --board <board-id>

# Create a new label
trello-cli label create --board <board-id> --name "Important" --color "red"

# Add label to card
trello-cli label add <card-id> <label-id>
```

#### Checklists

```bash
# List all checklists on a card
trello-cli checklist list --card <card-id>

# Create a new checklist
trello-cli checklist create --card <card-id> "Task List"

# Add item to checklist
trello-cli checklist add-item <checklist-id> "Task Item"
```

#### Members

```bash
# Get member information
trello-cli member get <username-or-id>

# List member's boards
trello-cli member boards <username-or-id>
```

#### Attachments

```bash
# List all attachments on a card
trello-cli attachment list --card <card-id>

# Add attachment to card
trello-cli attachment add --card <card-id> <url>
```

### Output Formats

#### JSON (Default)

```bash
trello-cli board list
# Outputs structured JSON
```

#### Markdown

```bash
trello-cli board list --format markdown
# Outputs formatted Markdown tables and sections
```

### LLM Optimization Features

#### Field Filtering

```bash
# Only include specific fields
trello-cli card list --list <list-id> --fields name,desc,due

# Verbose output with all fields
trello-cli card get <card-id> --verbose
```

#### Token Limits

```bash
# Limit output to 2000 tokens
trello-cli board list --max-tokens 2000
```

#### Quiet Mode

```bash
# Minimal output for scripting
trello-cli card create --list <list-id> "New Card" --quiet
```

### Batch Operations

#### From File

Create a batch file `operations.json`:

```json
{
  "operations": [
    {
      "type": "card",
      "resource": "card",
      "action": "create",
      "data": {
        "name": "Task 1",
        "list_id": "list-id-1"
      }
    },
    {
      "type": "card",
      "resource": "card", 
      "action": "create",
      "data": {
        "name": "Task 2",
        "list_id": "list-id-2"
      }
    }
  ],
  "continue_on_error": true
}
```

Execute batch operations:

```bash
trello-cli batch file operations.json
```

#### From Stdin

```bash
cat operations.json | trello-cli batch stdin
```

### Configuration Management

```bash
# Show current configuration
trello-cli config show

# Set configuration values
trello-cli config set --api-key "key" --token "token" --default-format json

# Show config file path
trello-cli config path
```

### Schema Output (for LLM Integration)

```bash
# Get complete CLI schema in JSON format
trello-cli schema

# Save schema to file
trello-cli schema > trello-schema.json

# Query schema with jq
trello-cli schema | jq '.subcommands[] | select(.name | startswith("board"))'
```

The schema command outputs a comprehensive JSON schema of all commands, arguments, flags, and usage patterns - perfect for LLM consumption and programmatic discovery.

## LLM Integration Examples

### Getting Board Context for LLM

```bash
# Get board summary with key information
trello-cli board get <board-id> --fields name,desc,url --format json

# Get all cards with essential fields
trello-cli card list --list <list-id> --fields name,desc,labels,due --format json --max-tokens 3000
```

### Batch Processing for LLM Workflows

```bash
# Process multiple operations from LLM-generated JSON
echo '{"operations":[{"type":"card","resource":"card","action":"create","data":{"name":"LLM Generated Task","list_id":"list-id"}}]}' | trello-cli batch stdin --format json
```

### Scripting Integration

```bash
#!/bin/bash
# Create cards from a list
while IFS= read -r task; do
    trello-cli card create --list "$LIST_ID" "$task" --quiet
done < tasks.txt
```

## Configuration Reference

### Config File Format

```yaml
api_key: your-trello-api-key
token: your-trello-token
default_format: json  # or markdown
max_tokens: 4000      # 0 = unlimited
```

### Environment Variables

- `TRELLO_API_KEY`: Your Trello API key
- `TRELLO_TOKEN`: Your Trello access token

### Global Flags

- `--api-key`: Override API key
- `--token`: Override token
- `--format, -f`: Output format (markdown, json)
- `--fields`: Comma-separated list of fields to include
- `--max-tokens`: Maximum tokens in output (0 = unlimited)
- `--verbose, -v`: Verbose output
- `--quiet, -q`: Quiet mode (minimal output)
- `--debug`: Debug mode (show API calls)

## Error Handling

The CLI provides detailed error messages and appropriate exit codes:

- `0`: Success
- `1`: General error
- `2`: Authentication error
- `3`: API error

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:

1. Check the [GitHub Issues](https://github.com/danbruder/trello-cli/issues)
2. Create a new issue with detailed information
3. Include your Go version, OS, and error messages

