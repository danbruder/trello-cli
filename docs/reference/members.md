# Members

View member information and manage member-related operations.

## Commands

### `get`
Get detailed information about a member.

```bash
trello-cli member get <username-or-id> [flags]
```

**Arguments:**
- `<username-or-id>` - The username or ID of the member to retrieve

**Examples:**
```bash
# Get member details by username
trello-cli member get john_doe

# Get member details by ID
trello-cli member get 5f8b8c8d8e8f8a8b8c8d8e8f

# Get member details in JSON format
trello-cli member get john_doe --format json

# Get only specific fields
trello-cli member get john_doe --fields username,fullName,avatarHash
```

### `boards`
List all boards that a member has access to.

```bash
trello-cli member boards <username-or-id> [flags]
```

**Arguments:**
- `<username-or-id>` - The username or ID of the member

**Examples:**
```bash
# List member's boards
trello-cli member boards john_doe

# List member's boards in JSON format
trello-cli member boards john_doe --format json

# List boards with specific fields
trello-cli member boards john_doe --fields name,desc,closed
```

## Common Use Cases

### Team Management
```bash
# Get team member information
trello-cli member get john_doe --fields username,fullName,avatarHash

# List all boards a member has access to
trello-cli member boards john_doe --fields name,desc,closed
```

### LLM Integration
```bash
# Get member context for LLM processing
trello-cli member get john_doe --fields username,fullName --format json

# Get member's board access for context
trello-cli member boards john_doe --fields name,closed --format json --max-tokens 2000
```

### Team Onboarding
```bash
#!/bin/bash
# Check team member access
MEMBER="new_team_member"
echo "Checking access for $MEMBER..."

# Get member info
trello-cli member get "$MEMBER" --fields username,fullName

# List their boards
echo "Boards accessible to $MEMBER:"
trello-cli member boards "$MEMBER" --fields name,desc
```
