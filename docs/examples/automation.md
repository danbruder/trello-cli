# Automation Examples

Learn how to automate Trello workflows using shell scripts and external tools.

## Shell Scripting Automation

### Basic Task Automation

```bash
#!/bin/bash
# Automated task creation from text file

TASKS_FILE="daily-tasks.txt"
LIST_ID="your-list-id"

# Create cards from task list
while IFS= read -r task; do
    if [ -n "$task" ]; then
        trlo card create --list "$LIST_ID" "$task" --quiet
        echo "Created card: $task"
    fi
done < "$TASKS_FILE"
```

### Scheduled Task Processing

```bash
#!/bin/bash
# Cron job to process daily tasks

# Set up environment
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"

# Process tasks from different sources
process_email_tasks() {
    # Process tasks from email (example)
    echo "Processing email tasks..."
    # Implementation depends on email processing setup
}

process_calendar_events() {
    # Process calendar events as tasks
    echo "Processing calendar events..."
    # Implementation depends on calendar integration
}

# Run daily processing
process_email_tasks
process_calendar_events
```

## File System Automation

### Directory Monitoring

```bash
#!/bin/bash
# Monitor directory for new files and create Trello cards

WATCH_DIR="/path/to/watch"
LIST_ID="your-list-id"
BOARD_ID="your-board-id"

# Function to create card for new file
create_file_card() {
    local file_path="$1"
    local file_name=$(basename "$file_path")
    local file_size=$(du -h "$file_path" | cut -f1)
    
    trlo card create --list "$LIST_ID" "New file: $file_name" \
        --desc "File: $file_path\nSize: $file_size\nAdded: $(date)" \
        --quiet
}

# Watch for file creation
inotifywait -m -e create "$WATCH_DIR" | while read -r path action file; do
    create_file_card "$path$file"
done
```

### Log File Processing

```bash
#!/bin/bash
# Process application logs and create Trello cards for issues

LOG_FILE="/var/log/application.log"
ERROR_LIST_ID="error-list-id"
WARNING_LIST_ID="warning-list-id"

# Process errors
grep "ERROR" "$LOG_FILE" | tail -5 | while read -r line; do
    timestamp=$(echo "$line" | cut -d' ' -f1-2)
    error_msg=$(echo "$line" | cut -d' ' -f3-)
    
    trlo card create --list "$ERROR_LIST_ID" \
        "Error: $timestamp" \
        --desc "$error_msg" \
        --quiet
done

# Process warnings
grep "WARNING" "$LOG_FILE" | tail -3 | while read -r line; do
    timestamp=$(echo "$line" | cut -d' ' -f1-2)
    warning_msg=$(echo "$line" | cut -d' ' -f3-)
    
    trlo card create --list "$WARNING_LIST_ID" \
        "Warning: $timestamp" \
        --desc "$warning_msg" \
        --quiet
done
```

## Database Automation

### MySQL Task Sync

```bash
#!/bin/bash
# Sync MySQL database tasks with Trello

DB_HOST="localhost"
DB_USER="user"
DB_PASS="password"
DB_NAME="taskdb"
LIST_ID="your-list-id"

# Query database for new tasks
mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
    SELECT id, title, description, priority 
    FROM tasks 
    WHERE trello_card_id IS NULL 
    AND status = 'pending'
" | while read -r id title desc priority; do
    if [ "$id" != "id" ]; then  # Skip header row
        # Create Trello card
        CARD_ID=$(trlo card create --list "$LIST_ID" "$title" --desc "$desc" --quiet)
        
        # Update database with Trello card ID
        mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
            UPDATE tasks 
            SET trello_card_id = '$CARD_ID' 
            WHERE id = '$id'
        "
        
        echo "Synced task: $title"
    fi
done
```

### PostgreSQL Integration

```bash
#!/bin/bash
# PostgreSQL to Trello automation

DB_CONNECTION="postgresql://user:password@localhost/database"
LIST_ID="your-list-id"

# Create cards from PostgreSQL query
psql "$DB_CONNECTION" -t -c "
    SELECT title, description, due_date 
    FROM projects 
    WHERE status = 'active' 
    AND trello_card_id IS NULL
" | while IFS='|' read -r title desc due_date; do
    # Clean up whitespace
    title=$(echo "$title" | xargs)
    desc=$(echo "$desc" | xargs)
    due_date=$(echo "$due_date" | xargs)
    
    # Create card with due date
    CARD_ID=$(trlo card create --list "$LIST_ID" "$title" --desc "$desc" --quiet)
    
    # Update database
    psql "$DB_CONNECTION" -c "
        UPDATE projects 
        SET trello_card_id = '$CARD_ID' 
        WHERE title = '$title'
    "
done
```

## External Service Automation

### GitHub Issue Sync

```bash
#!/bin/bash
# Sync GitHub issues with Trello cards

GITHUB_TOKEN="your-github-token"
GITHUB_REPO="owner/repo"
LIST_ID="your-list-id"

# Fetch open issues
curl -s -H "Authorization: token $GITHUB_TOKEN" \
    "https://api.github.com/repos/$GITHUB_REPO/issues?state=open" | \
    jq -r '.[] | "\(.title)|\(.body)|\(.html_url)"' | while IFS='|' read -r title body url; do
    
    # Create Trello card
    trlo card create --list "$LIST_ID" "$title" \
        --desc "$body\n\nGitHub: $url" \
        --quiet
    
    echo "Synced issue: $title"
done
```

### Slack Integration

```bash
#!/bin/bash
# Process Slack messages and create Trello cards

SLACK_WEBHOOK_URL="your-slack-webhook-url"
LIST_ID="your-list-id"

# Function to process Slack message
process_slack_message() {
    local message="$1"
    local user="$2"
    local channel="$3"
    
    # Create card with Slack context
    trlo card create --list "$LIST_ID" \
        "Slack: $message" \
        --desc "From: $user\nChannel: $channel\nTime: $(date)" \
        --quiet
}

# Example: Process incoming webhook
if [ -n "$SLACK_PAYLOAD" ]; then
    MESSAGE=$(echo "$SLACK_PAYLOAD" | jq -r '.text')
    USER=$(echo "$SLACK_PAYLOAD" | jq -r '.user_name')
    CHANNEL=$(echo "$SLACK_PAYLOAD" | jq -r '.channel_name')
    
    process_slack_message "$MESSAGE" "$USER" "$CHANNEL"
fi
```

## Advanced Automation Patterns

### Conditional Workflows

```bash
#!/bin/bash
# Conditional automation based on various factors

check_and_create_tasks() {
    local condition="$1"
    local task_name="$2"
    local task_desc="$3"
    
    case "$condition" in
        "high_load")
            if [ $(loadavg | cut -d' ' -f1 | cut -d'.' -f1) -gt 5 ]; then
                trlo card create --list "$HIGH_PRIORITY_LIST" "$task_name" --desc "$task_desc"
            fi
            ;;
        "disk_space")
            if [ $(df / | tail -1 | awk '{print $5}' | sed 's/%//') -gt 80 ]; then
                trlo card create --list "$URGENT_LIST" "$task_name" --desc "$task_desc"
            fi
            ;;
        "memory_usage")
            if [ $(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100.0}') -gt 90 ]; then
                trlo card create --list "$SYSTEM_LIST" "$task_name" --desc "$task_desc"
            fi
            ;;
    esac
}

# Monitor system and create tasks
check_and_create_tasks "high_load" "High CPU Load" "System experiencing high CPU usage"
check_and_create_tasks "disk_space" "Low Disk Space" "Disk usage above 80%"
check_and_create_tasks "memory_usage" "High Memory Usage" "Memory usage above 90%"
```

### Error Handling and Retry Logic

```bash
#!/bin/bash
# Robust automation with error handling

MAX_RETRIES=3
RETRY_DELAY=5

create_card_with_retry() {
    local list_id="$1"
    local title="$2"
    local description="$3"
    local retry_count=0
    
    while [ $retry_count -lt $MAX_RETRIES ]; do
        if trlo card create --list "$list_id" "$title" --desc "$description" --quiet; then
            echo "✓ Card created: $title"
            return 0
        else
            echo "✗ Failed to create card (attempt $((retry_count + 1))/$MAX_RETRIES): $title"
            retry_count=$((retry_count + 1))
            sleep $RETRY_DELAY
        fi
    done
    
    echo "✗ Failed to create card after $MAX_RETRIES attempts: $title"
    return 1
}

# Use the robust function
create_card_with_retry "$LIST_ID" "Automated Task" "Created by automation script"
```

### Batch Processing with Progress Tracking

```bash
#!/bin/bash
# Process large batches with progress tracking

process_large_batch() {
    local input_file="$1"
    local list_id="$2"
    local total_lines=$(wc -l < "$input_file")
    local current_line=0
    
    echo "Processing $total_lines items..."
    
    while IFS= read -r item; do
        current_line=$((current_line + 1))
        
        if trlo card create --list "$list_id" "$item" --quiet; then
            echo "✓ [$current_line/$total_lines] Created: $item"
        else
            echo "✗ [$current_line/$total_lines] Failed: $item"
        fi
        
        # Progress indicator
        if [ $((current_line % 10)) -eq 0 ]; then
            echo "Progress: $current_line/$total_lines ($((current_line * 100 / total_lines))%)"
        fi
    done < "$input_file"
    
    echo "Batch processing complete!"
}

# Process large task list
process_large_batch "large-task-list.txt" "$LIST_ID"
```
