# Boards

Manage Trello boards including listing, creating, updating, and deleting boards.

## Commands

### `list`
List all boards you have access to.

```bash
trello-cli board list [flags]
```

**Examples:**
```bash
# List all boards
trello-cli board list

# List boards in JSON format
trello-cli board list --format json

# List boards with specific fields
trello-cli board list --fields name,desc,url
```

### `get`
Get detailed information about a specific board.

```bash
trello-cli board get <board-id> [flags]
```

**Arguments:**
- `<board-id>` - The ID of the board to retrieve

**Examples:**
```bash
# Get board details
trello-cli board get 5f8b8c8d8e8f8a8b8c8d8e8f

# Get board details in JSON format
trello-cli board get 5f8b8c8d8e8f8a8b8c8d8e8f --format json

# Get only specific fields
trello-cli board get 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,desc,closed
```

### `create`
Create a new board.

```bash
trello-cli board create <name> [flags]
```

**Arguments:**
- `<name>` - The name of the board to create

**Examples:**
```bash
# Create a new board
trello-cli board create "My New Board"

# Create board with description
trello-cli board create "Project Board" --desc "Board for project management"
```

### `add-member`
Add a member to a board.

```bash
trello-cli board add-member <board-id> <email> [flags]
```

**Arguments:**
- `<board-id>` - The ID of the board
- `<email>` - Email address of the member to add

**Examples:**
```bash
# Add member to board
trello-cli board add-member 5f8b8c8d8e8f8a8b8c8d8e8f user@example.com
```

### `delete`
Delete a board permanently.

```bash
trello-cli board delete <board-id> [flags]
```

**Arguments:**
- `<board-id>` - The ID of the board to delete

**Examples:**
```bash
# Delete a board
trello-cli board delete 5f8b8c8d8e8f8a8b8c8d8e8f
```

## Common Use Cases

### Get Board Context for LLM
```bash
# Get board summary with key information
trello-cli board get <board-id> --fields name,desc,url --format json

# Get all boards with essential fields
trello-cli board list --fields name,desc,closed --format json --max-tokens 3000
```

### Board Management Workflow
```bash
# 1. List boards to find the one you want
trello-cli board list

# 2. Get board details
trello-cli board get <board-id>

# 3. Add team members
trello-cli board add-member <board-id> team@company.com

# 4. List lists on the board
trello-cli list list --board <board-id>
```
