# Cards

Manage Trello cards including listing, creating, updating, moving, copying, and deleting cards.

## Commands

### `list`
List all cards in a specific list.

```bash
trello-cli card list --list <list-id> [flags]
```

**Flags:**
- `--list` - The ID of the list to list cards from

**Examples:**
```bash
# List cards in a list
trello-cli card list --list 5f8b8c8d8e8f8a8b8c8d8e8f

# List cards with specific fields
trello-cli card list --list 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,desc,due

# List cards in JSON format
trello-cli card list --list 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `get`
Get detailed information about a specific card.

```bash
trello-cli card get <card-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to retrieve

**Examples:**
```bash
# Get card details
trello-cli card get 5f8b8c8d8e8f8a8b8c8d8e8f

# Get card details in JSON format
trello-cli card get 5f8b8c8d8e8f8a8b8c8d8e8f --format json

# Get only specific fields
trello-cli card get 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,desc,labels,due
```

### `create`
Create a new card in a list.

```bash
trello-cli card create --list <list-id> <name> [flags]
```

**Arguments:**
- `<name>` - The name of the card to create

**Flags:**
- `--list` - The ID of the list to create the card in

**Examples:**
```bash
# Create a new card
trello-cli card create --list 5f8b8c8d8e8f8a8b8c8d8e8f "My New Card"

# Create card with description
trello-cli card create --list 5f8b8c8d8e8f8a8b8c8d8e8f "Task Card" --desc "Description of the task"

# Create card quietly for scripting
trello-cli card create --list 5f8b8c8d8e8f8a8b8c8d8e8f "New Task" --quiet
```

### `move`
Move a card to another list.

```bash
trello-cli card move <card-id> --list <target-list-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to move

**Flags:**
- `--list` - The ID of the target list

**Examples:**
```bash
# Move card to another list
trello-cli card move 5f8b8c8d8e8f8a8b8c8d8e8f --list 5f8b8c8d8e8f8a8b8c8d8e8g
```

### `copy`
Copy a card to another list.

```bash
trello-cli card copy <card-id> --list <target-list-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to copy

**Flags:**
- `--list` - The ID of the target list

**Examples:**
```bash
# Copy card to another list
trello-cli card copy 5f8b8c8d8e8f8a8b8c8d8e8f --list 5f8b8c8d8e8f8a8b8c8d8e8g
```

### `archive`
Archive a card (soft delete).

```bash
trello-cli card archive <card-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to archive

**Examples:**
```bash
# Archive a card
trello-cli card archive 5f8b8c8d8e8f8a8b8c8d8e8f
```

### `delete`
Permanently delete a card.

```bash
trello-cli card delete <card-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to delete

**Examples:**
```bash
# Delete a card permanently
trello-cli card delete 5f8b8c8d8e8f8a8b8c8d8e8f
```

## Common Use Cases

### Card Management Workflow
```bash
# 1. List cards in a list
trello-cli card list --list <list-id>

# 2. Create a new card
trello-cli card create --list <list-id> "New Task"

# 3. Move card to done list
trello-cli card move <card-id> --list <done-list-id>

# 4. Archive completed cards
trello-cli card archive <card-id>
```

### LLM Integration
```bash
# Get cards with essential fields for LLM processing
trello-cli card list --list <list-id> --fields name,desc,labels,due --format json --max-tokens 3000

# Get specific card details for context
trello-cli card get <card-id> --fields name,desc,labels,attachments --format json
```

### Batch Card Operations
```bash
# Create multiple cards from a list
while IFS= read -r task; do
    trello-cli card create --list "$LIST_ID" "$task" --quiet
done < tasks.txt
```
