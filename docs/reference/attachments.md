# Attachments

Manage file attachments on Trello cards.

## Commands

### `list`
List all attachments on a card.

```bash
trello-cli attachment list --card <card-id> [flags]
```

**Flags:**
- `--card` - The ID of the card to list attachments from

**Examples:**
```bash
# List attachments on a card
trello-cli attachment list --card 5f8b8c8d8e8f8a8b8c8d8e8f

# List attachments with specific fields
trello-cli attachment list --card 5f8b8c8d8e8f8a8b8c8d8e8f --fields name,url,mimeType

# List attachments in JSON format
trello-cli attachment list --card 5f8b8c8d8e8f8a8b8c8d8e8f --format json
```

### `add`
Add an attachment to a card.

```bash
trello-cli attachment add --card <card-id> <url> [flags]
```

**Arguments:**
- `<url>` - The URL of the file to attach

**Flags:**
- `--card` - The ID of the card to add the attachment to

**Examples:**
```bash
# Add an attachment from URL
trello-cli attachment add --card 5f8b8c8d8e8f8a8b8c8d8e8f "https://example.com/file.pdf"

# Add attachment quietly for scripting
trello-cli attachment add --card 5f8b8c8d8e8f8a8b8c8d8e8f "https://example.com/image.png" --quiet
```

## Common Use Cases

### Document Management
```bash
# List all attachments on a card
trello-cli attachment list --card <card-id>

# Add documentation to a card
trello-cli attachment add --card <card-id> "https://docs.example.com/api-reference.pdf"

# Add images to a card
trello-cli attachment add --card <card-id> "https://example.com/screenshot.png"
```

### LLM Integration
```bash
# Get attachment information for context
trello-cli attachment list --card <card-id> --fields name,url,mimeType --format json

# Get card with attachment details
trello-cli card get <card-id> --fields name,attachments --format json
```

### Automation Scripts
```bash
#!/bin/bash
# Add multiple attachments to a card
CARD_ID="your-card-id"
ATTACHMENTS=(
    "https://example.com/doc1.pdf"
    "https://example.com/doc2.docx"
    "https://example.com/image1.png"
)

for attachment in "${ATTACHMENTS[@]}"; do
    trello-cli attachment add --card "$CARD_ID" "$attachment" --quiet
done
```
