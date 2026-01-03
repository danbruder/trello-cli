# Schema

Output a complete CLI schema in JSON format for LLM consumption.

## Command

### `schema`
Output comprehensive JSON schema of all commands, subcommands, flags, and arguments.

```bash
trello-cli schema [flags]
```

**Examples:**
```bash
# Get complete schema
trello-cli schema

# Save schema to file
trello-cli schema > trello-schema.json

# Parse with jq
trello-cli schema | jq '.subcommands[] | select(.name | startswith("board"))'
```

## Schema Structure

The schema command outputs a JSON object with the following structure:

```json
{
  "name": "trello-cli",
  "description": "CLI description",
  "usage": "trello-cli [command]",
  "global_flags": [
    {
      "name": "flag-name",
      "short": "f",
      "description": "Flag description",
      "type": "string|int|bool|[]string",
      "default": "default-value",
      "required": true|false
    }
  ],
  "subcommands": [
    {
      "name": "command-name",
      "description": "Command description",
      "usage": "trello-cli command [args] [flags]",
      "arguments": [
        {
          "name": "arg-name",
          "description": "Argument description",
          "required": true|false,
          "type": "string"
        }
      ],
      "flags": [
        {
          "name": "flag-name",
          "description": "Flag description",
          "type": "string",
          "required": true|false
        }
      ],
      "examples": [
        "trello-cli command example"
      ]
    }
  ],
  "examples": [
    "trello-cli command example"
  ]
}
```

## Use Cases

### LLM Integration

The schema command is specifically designed for LLM consumption, providing:

- **Complete API Surface**: All 31 commands in one response
- **Structured Metadata**: Arguments, flags, types, requirements
- **Usage Patterns**: Clear syntax for each command
- **Examples**: Practical usage demonstrations

```bash
# LLM can query the schema to understand available operations
trello-cli schema
```

### Documentation Generation

```bash
# Extract all board commands
trello-cli schema | jq '.subcommands[] | select(.name | startswith("board"))'

# List all commands with their descriptions
trello-cli schema | jq '.subcommands[] | {name, description}'

# Get all required arguments
trello-cli schema | jq '.subcommands[] | select(.arguments != null) | {name, arguments: .arguments[] | select(.required == true)}'
```

### API Discovery

```bash
# Find commands that accept a specific flag
trello-cli schema | jq '.subcommands[] | select(.flags != null) | select(.flags[] | .name == "list")'

# Get all commands with examples
trello-cli schema | jq '.subcommands[] | select(.examples != null) | {name, examples}'
```

### Programmatic Usage

```python
import subprocess
import json

# Get schema
result = subprocess.run(['trello-cli', 'schema'], capture_output=True, text=True)
schema = json.loads(result.stdout)

# Discover available commands
for cmd in schema['subcommands']:
    print(f"{cmd['name']}: {cmd['description']}")

# Find commands that require a board-id
board_commands = [
    cmd for cmd in schema['subcommands']
    if cmd.get('arguments') and any(arg['name'] == 'board-id' for arg in cmd['arguments'])
]
```

```javascript
const { execSync } = require('child_process');

// Get schema
const schemaOutput = execSync('trello-cli schema').toString();
const schema = JSON.parse(schemaOutput);

// Build command dynamically
const createCardCmd = schema.subcommands.find(cmd => cmd.name === 'card create');
console.log('Usage:', createCardCmd.usage);
console.log('Required arguments:', createCardCmd.arguments.filter(arg => arg.required));
console.log('Available flags:', createCardCmd.flags);
```

## Schema Contents

The schema includes documentation for:

### Board Commands (5)
- board list
- board get
- board create (with --desc flag)
- board delete
- board add-member

### Card Commands (7)
- card list
- card get
- card create (with --desc flag)
- card move
- card copy
- card archive
- card delete

### List Commands (4)
- list list
- list get
- list create
- list archive

### Label Commands (3)
- label list
- label create
- label add

### Checklist Commands (3)
- checklist list
- checklist create
- checklist add-item

### Member Commands (2)
- member get
- member boards

### Attachment Commands (2)
- attachment list
- attachment add

### Batch Commands (2)
- batch file
- batch stdin

### Config Commands (3)
- config show
- config set
- config path

### Global Flags (8)
- --api-key
- --token
- --format / -f
- --fields
- --max-tokens
- --verbose / -v
- --quiet / -q
- --debug

## Benefits for LLM Workflows

1. **Single Source of Truth**: One command provides complete API documentation
2. **Type Safety**: Know exactly what types each argument/flag expects
3. **Required vs Optional**: Understand which parameters are mandatory
4. **Examples**: Learn from practical usage demonstrations
5. **Discoverability**: Programmatically explore available operations
6. **Validation**: Validate command construction before execution

## Example Output

```bash
trello-cli schema | jq '.subcommands[0]'
```

```json
{
  "name": "board list",
  "description": "List all boards accessible to the authenticated user",
  "usage": "trello-cli board list [flags]",
  "examples": [
    "trello-cli board list",
    "trello-cli board list --format json",
    "trello-cli board list --fields name,desc,url"
  ]
}
```
