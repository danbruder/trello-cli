# Lists

Manage Trello lists including listing, creating, updating, and archiving lists.

## Commands

### `list`
List all lists on a board.

```bash
trlo list list --board <board-id> [flags]
```

**Flags:**
- `--board` - The ID of the board to list lists from

**Examples:**
```bash
# List lists on a board
trlo list list --board 5f8b8c8d8e8f8a8b8c8d8e8f

# List lists with specific fields
trlo list list --board 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,closed

# List lists in JSON format
trlo list list --board 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `get`
Get detailed information about a specific list.

```bash
trlo list get <list-id> [flags]
```

**Arguments:**
- `<list-id>` - The ID of the list to retrieve

**Examples:**
```bash
# Get list details
trlo list get 5f8b8c8d8e8f8a8b8c8d8e8f

# Get list details in JSON format
trlo list get 5f8b8c8d8e8f8a8b8c8d8e8f --format json

# Get only specific fields
trlo list get 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,closed,pos
```

### `create`
Create a new list on a board.

```bash
trlo list create --board <board-id> <name> [flags]
```

**Arguments:**
- `<name>` - The name of the list to create

**Flags:**
- `--board` - The ID of the board to create the list on

**Examples:**
```bash
# Create a new list
trlo list create --board 5f8b8c8d8e8f8a8b8c8d8e8f "New List"

# Create list quietly for scripting
trlo list create --board 5f8b8c8d8e8f8a8b8c8d8e8f "Backlog" --quiet
```

### `archive`
Archive a list (soft delete).

```bash
trlo list archive <list-id> [flags]
```

**Arguments:**
- `<list-id>` - The ID of the list to archive

**Examples:**
```bash
# Archive a list
trlo list archive 5f8b8c8d8e8f8a8b8c8d8e8f
```

## Common Use Cases

### Board Setup Workflow
```bash
# 1. List existing lists on a board
trlo list list --board <board-id>

# 2. Create new lists for project phases
trlo list create --board <board-id> "To Do"
trlo list create --board <board-id> "In Progress"
trlo list create --board <board-id> "Done"

# 3. Archive old lists
trlo list archive <old-list-id>
```

### LLM Integration
```bash
# Get all lists with essential information
trlo list list --board <board-id> --fields name,closed --format json

# Get specific list details for context
trlo list get <list-id> --fields name,closed,pos --format json
```

### Automation Scripts
```bash
#!/bin/bash
# Create standard lists for a new project
BOARD_ID="your-board-id"
trlo list create --board "$BOARD_ID" "Backlog" --quiet
trlo list create --board "$BOARD_ID" "Sprint" --quiet
trlo list create --board "$BOARD_ID" "In Progress" --quiet
trlo list create --board "$BOARD_ID" "Review" --quiet
trlo list create --board "$BOARD_ID" "Done" --quiet
```
