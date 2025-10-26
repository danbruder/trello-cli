# API Integration Examples

Learn how to integrate Trello CLI with external APIs and services for powerful automation workflows.

## REST API Integration

### Basic HTTP Integration

```bash
#!/bin/bash
# Create cards from external API data

# Fetch tasks from external API
TASKS=$(curl -s "https://api.example.com/tasks" | jq -r '.[].title')

# Create cards for each task
while IFS= read -r task; do
    trlo card create --list "$LIST_ID" "$task" --quiet
done <<< "$TASKS"
```

### Webhook Integration

```bash
#!/bin/bash
# Process webhook payload and create Trello cards

# Parse webhook JSON
TITLE=$(echo "$WEBHOOK_PAYLOAD" | jq -r '.title')
DESCRIPTION=$(echo "$WEBHOOK_PAYLOAD" | jq -r '.description')
PRIORITY=$(echo "$WEBHOOK_PAYLOAD" | jq -r '.priority')

# Create card based on priority
if [ "$PRIORITY" = "high" ]; then
    LIST_ID="$HIGH_PRIORITY_LIST"
else
    LIST_ID="$NORMAL_LIST"
fi

trlo card create --list "$LIST_ID" "$TITLE" --desc "$DESCRIPTION"
```

## Database Integration

### MySQL Integration

```bash
#!/bin/bash
# Sync database records with Trello cards

# Query database for pending tasks
mysql -u user -p database -e "SELECT id, title, description FROM tasks WHERE status='pending'" | while read -r id title desc; do
    # Create Trello card
    CARD_ID=$(trlo card create --list "$LIST_ID" "$title" --desc "$desc" --quiet)
    
    # Update database with Trello card ID
    mysql -u user -p database -e "UPDATE tasks SET trello_card_id='$CARD_ID' WHERE id='$id'"
done
```

### PostgreSQL Integration

```bash
#!/bin/bash
# Create Trello cards from PostgreSQL query results

psql -d database -t -c "SELECT title, description FROM projects WHERE status='active'" | while IFS='|' read -r title desc; do
    # Clean up whitespace
    title=$(echo "$title" | xargs)
    desc=$(echo "$desc" | xargs)
    
    # Create card
    trlo card create --list "$LIST_ID" "$title" --desc "$desc"
done
```

## File System Integration

### Directory Monitoring

```bash
#!/bin/bash
# Monitor directory for new files and create Trello cards

WATCH_DIR="/path/to/watch"
LIST_ID="your-list-id"

# Watch for new files
inotifywait -m -e create "$WATCH_DIR" | while read -r path action file; do
    # Create card for new file
    trlo card create --list "$LIST_ID" "New file: $file" --desc "File created in $path"
done
```

### Log File Processing

```bash
#!/bin/bash
# Process log files and create Trello cards for errors

LOG_FILE="/var/log/application.log"
ERROR_LIST_ID="error-list-id"

# Extract errors from log
grep "ERROR" "$LOG_FILE" | tail -10 | while read -r line; do
    # Create card for each error
    trlo card create --list "$ERROR_LIST_ID" "Error: $(echo "$line" | cut -d' ' -f1-3)" --desc "$line"
done
```

## External Service Integration

### GitHub Integration

```bash
#!/bin/bash
# Create Trello cards from GitHub issues

# Fetch open issues
curl -s -H "Authorization: token $GITHUB_TOKEN" \
    "https://api.github.com/repos/owner/repo/issues?state=open" | \
    jq -r '.[] | "\(.title)|\(.body)"' | while IFS='|' read -r title body; do
    
    # Create Trello card
    trlo card create --list "$LIST_ID" "$title" --desc "$body"
done
```

### Slack Integration

```bash
#!/bin/bash
# Process Slack messages and create Trello cards

# Parse Slack webhook
MESSAGE=$(echo "$SLACK_PAYLOAD" | jq -r '.text')
USER=$(echo "$SLACK_PAYLOAD" | jq -r '.user_name')

# Create card with Slack context
trlo card create --list "$LIST_ID" "Slack: $MESSAGE" --desc "From: $USER"
```

## Batch Processing Integration

### CSV File Processing

```bash
#!/bin/bash
# Process CSV file and create Trello cards

CSV_FILE="tasks.csv"
LIST_ID="your-list-id"

# Convert CSV to batch operations
echo '{"operations":[' > batch.json

tail -n +2 "$CSV_FILE" | while IFS=',' read -r title description priority; do
    echo "    {\"type\":\"card\",\"resource\":\"card\",\"action\":\"create\",\"data\":{\"name\":\"$title\",\"desc\":\"$description\",\"list_id\":\"$LIST_ID\"}}," >> batch.json
done

# Remove trailing comma and close JSON
sed -i '$ s/,$//' batch.json
echo ']}' >> batch.json

# Execute batch operations
trlo batch file batch.json
```

### JSON API Response Processing

```bash
#!/bin/bash
# Process JSON API response and create batch operations

# Fetch data from API
API_RESPONSE=$(curl -s "https://api.example.com/tasks")

# Convert to batch operations
echo "$API_RESPONSE" | jq '{
  operations: [.tasks[] | {
    type: "card",
    resource: "card", 
    action: "create",
    data: {
      name: .title,
      desc: .description,
      list_id: "'$LIST_ID'"
    }
  }],
  continue_on_error: true
}' > batch.json

# Execute batch operations
trlo batch file batch.json
```

## Error Handling and Retry Logic

```bash
#!/bin/bash
# Robust API integration with retry logic

MAX_RETRIES=3
RETRY_DELAY=5

create_card_with_retry() {
    local list_id="$1"
    local title="$2"
    local description="$3"
    local retry_count=0
    
    while [ $retry_count -lt $MAX_RETRIES ]; do
        if trlo card create --list "$list_id" "$title" --desc "$description" --quiet; then
            echo "Card created successfully: $title"
            return 0
        else
            echo "Failed to create card (attempt $((retry_count + 1))/$MAX_RETRIES)"
            retry_count=$((retry_count + 1))
            sleep $RETRY_DELAY
        fi
    done
    
    echo "Failed to create card after $MAX_RETRIES attempts: $title"
    return 1
}

# Use the function
create_card_with_retry "$LIST_ID" "API Integration Task" "Created from external API"
```
