# Batch Operations

Execute multiple Trello operations from a file or stdin for automation and scripting.

## Commands

### `file`
Execute batch operations from a JSON file.

```bash
trello-cli batch file <file-path> [flags]
```

**Arguments:**
- `<file-path>` - Path to the JSON file containing batch operations

**Examples:**
```bash
# Execute operations from a file
trello-cli batch file operations.json

# Execute with specific format
trello-cli batch file operations.json --format json

# Execute quietly
trello-cli batch file operations.json --quiet
```

### `stdin`
Execute batch operations from stdin.

```bash
trello-cli batch stdin [flags]
```

**Examples:**
```bash
# Execute operations from stdin
cat operations.json | trello-cli batch stdin

# Execute with pipe
echo '{"operations":[...]}' | trello-cli batch stdin

# Execute with format
cat operations.json | trello-cli batch stdin --format json
```

## Batch Operation Format

The batch operation file should be a JSON file with the following structure:

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

### Operation Structure

Each operation in the `operations` array should have:

- `type`: The resource type (card, list, board, etc.)
- `resource`: The specific resource (card, list, board, etc.)
- `action`: The action to perform (create, update, delete, etc.)
- `data`: The data for the operation

### Batch Options

- `continue_on_error`: Whether to continue processing if an operation fails (default: false)

## Common Use Cases

### Project Setup
```json
{
  "operations": [
    {
      "type": "list",
      "resource": "list",
      "action": "create",
      "data": {
        "name": "Backlog",
        "board_id": "board-id"
      }
    },
    {
      "type": "list",
      "resource": "list",
      "action": "create",
      "data": {
        "name": "In Progress",
        "board_id": "board-id"
      }
    },
    {
      "type": "list",
      "resource": "list",
      "action": "create",
      "data": {
        "name": "Done",
        "board_id": "board-id"
      }
    }
  ],
  "continue_on_error": true
}
```

### Task Creation
```json
{
  "operations": [
    {
      "type": "card",
      "resource": "card",
      "action": "create",
      "data": {
        "name": "Implement user authentication",
        "list_id": "list-id",
        "desc": "Add login and registration functionality"
      }
    },
    {
      "type": "card",
      "resource": "card",
      "action": "create",
      "data": {
        "name": "Write API tests",
        "list_id": "list-id",
        "desc": "Create comprehensive test suite"
      }
    }
  ]
}
```

### LLM-Generated Operations
```bash
# Process LLM-generated batch operations
echo '{"operations":[{"type":"card","resource":"card","action":"create","data":{"name":"LLM Generated Task","list_id":"list-id"}}]}' | trello-cli batch stdin --format json
```

### Automation Scripts
```bash
#!/bin/bash
# Generate batch operations from a task list
TASKS_FILE="tasks.txt"
BOARD_ID="your-board-id"
LIST_ID="your-list-id"

# Create batch operations JSON
echo '{"operations":[' > operations.json

while IFS= read -r task; do
    echo "    {\"type\":\"card\",\"resource\":\"card\",\"action\":\"create\",\"data\":{\"name\":\"$task\",\"list_id\":\"$LIST_ID\"}}," >> operations.json
done < "$TASKS_FILE"

# Remove trailing comma and close JSON
sed -i '$ s/,$//' operations.json
echo ']}' >> operations.json

# Execute batch operations
trello-cli batch file operations.json
```
