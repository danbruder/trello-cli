# Lists

Manage Trello lists including listing, creating, updating, and archiving lists.

## Commands

### `list`
List all lists on a board.

```bash
trello-cli list list --board <board-id> [flags]
```

**Flags:**
- `--board` - The ID of the board to list lists from

**Examples:**
```bash
# List lists on a board
trello-cli list list --board 5f8b8c8d8e8f8a8b8c8d8e8f

# List lists with specific fields
trello-cli list list --board 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,closed

# List lists in JSON format
trello-cli list list --board 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `get`
Get detailed information about a specific list.

```bash
trello-cli list get <list-id> [flags]
```

**Arguments:**
- `<list-id>` - The ID of the list to retrieve

**Examples:**
```bash
# Get list details
trello-cli list get 5f8b8c8d8e8f8a8b8c8d8e8f

# Get list details in JSON format
trello-cli list get 5f8b8c8d8e8f8a8b8c8d8e8f --format json

# Get only specific fields
trello-cli list get 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,closed,pos
```

### `create`
Create a new list on a board.

```bash
trello-cli list create --board <board-id> <name> [flags]
```

**Arguments:**
- `<name>` - The name of the list to create

**Flags:**
- `--board` - The ID of the board to create the list on

**Examples:**
```bash
# Create a new list
trello-cli list create --board 5f8b8c8d8e8f8a8b8c8d8e8f "New List"

# Create list quietly for scripting
trello-cli list create --board 5f8b8c8d8e8f8a8b8c8d8e8f "Backlog" --quiet
```

### `archive`
Archive a list (soft delete).

```bash
trello-cli list archive <list-id> [flags]
```

**Arguments:**
- `<list-id>` - The ID of the list to archive

**Examples:**
```bash
# Archive a list
trello-cli list archive 5f8b8c8d8e8f8a8b8c8d8e8f
```

## Common Use Cases

### Board Setup Workflow
```bash
# 1. List existing lists on a board
trello-cli list list --board <board-id>

# 2. Create new lists for project phases
trello-cli list create --board <board-id> "To Do"
trello-cli list create --board <board-id> "In Progress"
trello-cli list create --board <board-id> "Done"

# 3. Archive old lists
trello-cli list archive <old-list-id>
```

### LLM Integration
```bash
# Get all lists with essential information
trello-cli list list --board <board-id> --fields name,closed --format json

# Get specific list details for context
trello-cli list get <list-id> --fields name,closed,pos --format json
```

### Automation Scripts
```bash
#!/bin/bash
# Create standard lists for a new project
BOARD_ID="your-board-id"
trello-cli list create --board "$BOARD_ID" "Backlog" --quiet
trello-cli list create --board "$BOARD_ID" "Sprint" --quiet
trello-cli list create --board "$BOARD_ID" "In Progress" --quiet
trello-cli list create --board "$BOARD_ID" "Review" --quiet
trello-cli list create --board "$BOARD_ID" "Done" --quiet
```
