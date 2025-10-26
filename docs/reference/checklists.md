# Checklists

Manage Trello checklists including listing, creating, and adding items to checklists.

## Commands

### `list`
List all checklists on a card.

```bash
trlo checklist list --card <card-id> [flags]
```

**Flags:**
- `--card` - The ID of the card to list checklists from

**Examples:**
```bash
# List checklists on a card
trlo checklist list --card 5f8b8c8d8e8f8a8b8c8d8e8f

# List checklists with specific fields
trlo checklist list --card 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,checkItems

# List checklists in JSON format
trlo checklist list --card 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `create`
Create a new checklist on a card.

```bash
trlo checklist create --card <card-id> <name> [flags]
```

**Arguments:**
- `<name>` - The name of the checklist to create

**Flags:**
- `--card` - The ID of the card to create the checklist on

**Examples:**
```bash
# Create a new checklist
trlo checklist create --card 5f8b8c8d8e8f8a8b8c8d8e8f "Task List"

# Create checklist quietly for scripting
trlo checklist create --card 5f8b8c8d8e8f8a8b8c8d8e8f "Implementation Tasks" --quiet
```

### `add-item`
Add an item to a checklist.

```bash
trlo checklist add-item <checklist-id> <item-name> [flags]
```

**Arguments:**
- `<checklist-id>` - The ID of the checklist to add the item to
- `<item-name>` - The name of the item to add

**Examples:**
```bash
# Add an item to a checklist
trlo checklist add-item 5f8b8c8d8e8f8a8b8c8d8e8f "Review code"

# Add item quietly for scripting
trlo checklist add-item 5f8b8c8d8e8f8a8b8c8d8e8f "Write tests" --quiet
```

## Common Use Cases

### Task Breakdown Workflow
```bash
# 1. Create a checklist for a card
trlo checklist create --card <card-id> "Implementation Tasks"

# 2. Add items to the checklist
trlo checklist add-item <checklist-id> "Design API"
trlo checklist add-item <checklist-id> "Write tests"
trlo checklist add-item <checklist-id> "Code review"
trlo checklist add-item <checklist-id> "Deploy"

# 3. List checklist items
trlo checklist list --card <card-id>
```

### Project Management
```bash
# Create checklists for different phases
trlo checklist create --card <card-id> "Planning Phase"
trlo checklist create --card <card-id> "Development Phase"
trlo checklist create --card <card-id> "Testing Phase"
trlo checklist create --card <card-id> "Deployment Phase"
```

### LLM Integration
```bash
# Get checklist data for LLM processing
trlo checklist list --card <card-id> --fields name,checkItems --format json

# Get specific checklist details
trlo checklist list --card <card-id> --format json --max-tokens 2000
```

### Automation Scripts
```bash
#!/bin/bash
# Create a standard checklist for feature cards
CARD_ID="your-card-id"
CHECKLIST_ID=$(trlo checklist create --card "$CARD_ID" "Feature Checklist" --quiet)

# Add standard items
trlo checklist add-item "$CHECKLIST_ID" "Requirements analysis" --quiet
trlo checklist add-item "$CHECKLIST_ID" "Design review" --quiet
trlo checklist add-item "$CHECKLIST_ID" "Implementation" --quiet
trlo checklist add-item "$CHECKLIST_ID" "Testing" --quiet
trlo checklist add-item "$CHECKLIST_ID" "Documentation" --quiet
trlo checklist add-item "$CHECKLIST_ID" "Code review" --quiet
```
