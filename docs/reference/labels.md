# Labels

Manage Trello labels including listing, creating, and adding labels to cards.

## Commands

### `list`
List all labels on a board.

```bash
trlo label list --board <board-id> [flags]
```

**Flags:**
- `--board` - The ID of the board to list labels from

**Examples:**
```bash
# List labels on a board
trlo label list --board 5f8b8c8d8e8f8a8b8c8d8e8f

# List labels with specific fields
trlo label list --board 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,color

# List labels in JSON format
trlo label list --board 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `create`
Create a new label on a board.

```bash
trlo label create --board <board-id> --name <name> --color <color> [flags]
```

**Flags:**
- `--board` - The ID of the board to create the label on
- `--name` - The name of the label
- `--color` - The color of the label

**Available Colors:**
- `red`, `yellow`, `orange`, `lightgreen`, `green`, `lightblue`, `blue`, `purple`, `pink`, `black`

**Examples:**
```bash
# Create a red "Important" label
trlo label create --board 5f8b8c8d8e8f8a8b8c8d8e8f --name "Important" --color "red"

# Create a green "Done" label
trlo label create --board 5f8b8c8d8e8f8a8b8c8d8e8f --name "Done" --color "green"

# Create label quietly for scripting
trlo label create --board 5f8b8c8d8e8f8a8b8c8d8e8f --name "Bug" --color "red" --quiet
```

### `add`
Add a label to a card.

```bash
trlo label add <card-id> <label-id> [flags]
```

**Arguments:**
- `<card-id>` - The ID of the card to add the label to
- `<label-id>` - The ID of the label to add

**Examples:**
```bash
# Add a label to a card
trlo label add 5f8b8c8d8e8f8a8b8c8d8e8f 5f8b8c8d8e8f8a8b8c8d8e8g

# Add label quietly for scripting
trlo label add 5f8b8c8d8e8f8a8b8c8d8e8f 5f8b8c8d8e8f8a8b8c8d8e8g --quiet
```

## Common Use Cases

### Label Setup Workflow
```bash
# 1. List existing labels
trlo label list --board <board-id>

# 2. Create standard labels
trlo label create --board <board-id> --name "High Priority" --color "red"
trlo label create --board <board-id> --name "Medium Priority" --color "yellow"
trlo label create --board <board-id> --name "Low Priority" --color "green"
trlo label create --board <board-id> --name "Bug" --color "purple"
trlo label create --board <board-id> --name "Feature" --color "blue"
```

### Card Labeling Workflow
```bash
# 1. Get label IDs
trlo label list --board <board-id> --format json

# 2. Add labels to cards
trlo label add <card-id> <label-id>

# 3. Check card labels
trlo card get <card-id> --fields labels
```

### LLM Integration
```bash
# Get all labels for context
trlo label list --board <board-id> --fields name,color --format json

# Get cards with their labels
trlo card list --list <list-id> --fields name,labels --format json
```

### Automation Scripts
```bash
#!/bin/bash
# Create standard labels for a project
BOARD_ID="your-board-id"
trlo label create --board "$BOARD_ID" --name "Critical" --color "red" --quiet
trlo label create --board "$BOARD_ID" --name "High" --color "orange" --quiet
trlo label create --board "$BOARD_ID" --name "Medium" --color "yellow" --quiet
trlo label create --board "$BOARD_ID" --name "Low" --color "green" --quiet
```
