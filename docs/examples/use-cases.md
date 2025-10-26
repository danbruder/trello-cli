# Use Cases

Real-world scenarios and use cases for Trello CLI in various contexts.

## Project Management

### Agile Development Workflow

```bash
#!/bin/bash
# Complete agile workflow automation

BOARD_ID="your-project-board"
SPRINT_LIST="sprint-list-id"
BACKLOG_LIST="backlog-list-id"
DONE_LIST="done-list-id"

# Sprint Planning
plan_sprint() {
    echo "Planning sprint..."
    
    # Move high-priority items to sprint
    trello-cli card list --list "$BACKLOG_LIST" --fields name,desc --format json | \
        jq -r '.[] | select(.desc | contains("high priority")) | .id' | \
        while read -r card_id; do
            trello-cli card move "$card_id" --list "$SPRINT_LIST"
        done
}

# Daily Standup Data
daily_standup() {
    echo "Daily Standup Report:"
    echo "===================="
    
    # Yesterday's completed work
    echo "Completed Yesterday:"
    trello-cli card list --list "$DONE_LIST" --fields name --format json | \
        jq -r '.[] | .name' | head -5
    
    # Today's planned work
    echo -e "\nPlanned for Today:"
    trello-cli card list --list "$SPRINT_LIST" --fields name --format json | \
        jq -r '.[] | .name' | head -3
    
    # Blockers
    echo -e "\nBlockers:"
    trello-cli card list --list "$SPRINT_LIST" --fields name,desc --format json | \
        jq -r '.[] | select(.desc | contains("blocked")) | .name'
}

# Sprint Review
sprint_review() {
    echo "Sprint Review Report:"
    echo "===================="
    
    # Completed items
    echo "Completed Items:"
    trello-cli card list --list "$DONE_LIST" --fields name --format json | \
        jq -r '.[] | .name'
    
    # Incomplete items (move back to backlog)
    echo -e "\nMoving incomplete items to backlog..."
    trello-cli card list --list "$SPRINT_LIST" --fields name --format json | \
        jq -r '.[] | .id' | \
        while read -r card_id; do
            trello-cli card move "$card_id" --list "$BACKLOG_LIST"
        done
}
```

### Content Management Workflow

```bash
#!/bin/bash
# Content creation and publishing workflow

CONTENT_BOARD="content-board-id"
IDEAS_LIST="ideas-list-id"
WRITING_LIST="writing-list-id"
REVIEW_LIST="review-list-id"
PUBLISHED_LIST="published-list-id"

# Content Planning
plan_content() {
    local topic="$1"
    local deadline="$2"
    
    # Create content card
    CARD_ID=$(trello-cli card create --list "$IDEAS_LIST" \
        "Content: $topic" \
        --desc "Deadline: $deadline\nStatus: Planning" \
        --quiet)
    
    # Add content checklist
    trello-cli checklist create --card "$CARD_ID" "Content Checklist" --quiet
    CHECKLIST_ID=$(trello-cli checklist list --card "$CARD_ID" --format json | jq -r '.[0].id')
    
    # Add standard content tasks
    trello-cli checklist add-item "$CHECKLIST_ID" "Research topic" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Create outline" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Write first draft" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Review and edit" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Final review" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Publish" --quiet
}

# Content Pipeline Management
manage_content_pipeline() {
    echo "Content Pipeline Status:"
    echo "======================="
    
    # Ideas stage
    echo "Ideas:"
    trello-cli card list --list "$IDEAS_LIST" --fields name --format json | \
        jq -r '.[] | .name'
    
    # Writing stage
    echo -e "\nIn Writing:"
    trello-cli card list --list "$WRITING_LIST" --fields name --format json | \
        jq -r '.[] | .name'
    
    # Review stage
    echo -e "\nIn Review:"
    trello-cli card list --list "$REVIEW_LIST" --fields name --format json | \
        jq -r '.[] | .name'
    
    # Published
    echo -e "\nPublished:"
    trello-cli card list --list "$PUBLISHED_LIST" --fields name --format json | \
        jq -r '.[] | .name'
}
```

## Team Collaboration

### Team Onboarding Automation

```bash
#!/bin/bash
# Automated team member onboarding

onboard_team_member() {
    local member_name="$1"
    local member_email="$2"
    local role="$3"
    
    echo "Onboarding $member_name ($role)..."
    
    # Create onboarding board
    BOARD_NAME="$member_name - Onboarding"
    BOARD_ID=$(trello-cli board create "$BOARD_NAME" --quiet)
    
    # Add member to board
    trello-cli board add-member "$BOARD_ID" "$member_email"
    
    # Create onboarding lists
    trello-cli list create --board "$BOARD_ID" "Pre-Start" --quiet
    trello-cli list create --board "$BOARD_ID" "First Week" --quiet
    trello-cli list create --board "$BOARD_ID" "First Month" --quiet
    trello-cli list create --board "$BOARD_ID" "Completed" --quiet
    
    # Add role-specific tasks
    case "$role" in
        "developer")
            add_developer_tasks "$BOARD_ID"
            ;;
        "designer")
            add_designer_tasks "$BOARD_ID"
            ;;
        "manager")
            add_manager_tasks "$BOARD_ID"
            ;;
    esac
    
    echo "✓ Onboarding board created: $BOARD_NAME"
}

add_developer_tasks() {
    local board_id="$1"
    local first_week_list=$(trello-cli list list --board "$board_id" --format json | jq -r '.[] | select(.name == "First Week") | .id')
    
    trello-cli card create --list "$first_week_list" "Set up development environment" --quiet
    trello-cli card create --list "$first_week_list" "Review codebase documentation" --quiet
    trello-cli card create --list "$first_week_list" "Complete coding standards training" --quiet
    trello-cli card create --list "$first_week_list" "Pair programming session" --quiet
}
```

### Meeting Management

```bash
#!/bin/bash
# Meeting preparation and follow-up automation

MEETING_BOARD="meeting-board-id"
UPCOMING_LIST="upcoming-meetings-id"
FOLLOWUP_LIST="followup-actions-id"

# Prepare meeting agenda
prepare_meeting() {
    local meeting_name="$1"
    local date="$2"
    local attendees="$3"
    
    # Create meeting card
    CARD_ID=$(trello-cli card create --list "$UPCOMING_LIST" \
        "$meeting_name - $date" \
        --desc "Attendees: $attendees\nStatus: Preparing" \
        --quiet)
    
    # Add meeting checklist
    CHECKLIST_ID=$(trello-cli checklist create --card "$CARD_ID" "Meeting Prep" --quiet)
    
    trello-cli checklist add-item "$CHECKLIST_ID" "Send agenda" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Prepare materials" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Book meeting room" --quiet
    trello-cli checklist add-item "$CHECKLIST_ID" "Send calendar invite" --quiet
}

# Post-meeting follow-up
post_meeting_followup() {
    local meeting_card_id="$1"
    local action_items="$2"
    
    # Move meeting to completed
    trello-cli card move "$meeting_card_id" --list "$FOLLOWUP_LIST"
    
    # Create action items
    echo "$action_items" | while IFS='|' read -r action assignee due_date; do
        trello-cli card create --list "$FOLLOWUP_LIST" \
            "Action: $action" \
            --desc "Assignee: $assignee\nDue: $due_date\nFrom: Meeting" \
            --quiet
    done
}
```

## System Administration

### Infrastructure Monitoring

```bash
#!/bin/bash
# System monitoring and alerting via Trello

MONITORING_BOARD="monitoring-board-id"
ALERTS_LIST="alerts-list-id"
MAINTENANCE_LIST="maintenance-list-id"

# Monitor system resources
monitor_system() {
    # CPU monitoring
    CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
        trello-cli card create --list "$ALERTS_LIST" \
            "High CPU Usage Alert" \
            --desc "CPU usage: ${CPU_USAGE}%\nTime: $(date)\nServer: $(hostname)" \
            --quiet
    fi
    
    # Memory monitoring
    MEMORY_USAGE=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100.0}')
    if [ "$MEMORY_USAGE" -gt 85 ]; then
        trello-cli card create --list "$ALERTS_LIST" \
            "High Memory Usage Alert" \
            --desc "Memory usage: ${MEMORY_USAGE}%\nTime: $(date)\nServer: $(hostname)" \
            --quiet
    fi
    
    # Disk space monitoring
    DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
    if [ "$DISK_USAGE" -gt 80 ]; then
        trello-cli card create --list "$ALERTS_LIST" \
            "Low Disk Space Alert" \
            --desc "Disk usage: ${DISK_USAGE}%\nTime: $(date)\nServer: $(hostname)" \
            --quiet
    fi
}

# Schedule maintenance tasks
schedule_maintenance() {
    local task="$1"
    local scheduled_date="$2"
    
    trello-cli card create --list "$MAINTENANCE_LIST" \
        "Maintenance: $task" \
        --desc "Scheduled: $scheduled_date\nServer: $(hostname)\nStatus: Pending" \
        --quiet
}
```

### Backup Management

```bash
#!/bin/bash
# Backup monitoring and management

BACKUP_BOARD="backup-board-id"
DAILY_BACKUP_LIST="daily-backups-id"
WEEKLY_BACKUP_LIST="weekly-backups-id"
FAILED_BACKUP_LIST="failed-backups-id"

# Monitor backup status
check_backup_status() {
    local backup_type="$1"
    local backup_path="$2"
    
    if [ -f "$backup_path" ]; then
        local backup_size=$(du -h "$backup_path" | cut -f1)
        local backup_date=$(stat -c %y "$backup_path" | cut -d' ' -f1)
        
        trello-cli card create --list "$DAILY_BACKUP_LIST" \
            "$backup_type Backup - $backup_date" \
            --desc "Size: $backup_size\nPath: $backup_path\nStatus: Success" \
            --quiet
    else
        trello-cli card create --list "$FAILED_BACKUP_LIST" \
            "Failed: $backup_type Backup" \
            --desc "Path: $backup_path\nTime: $(date)\nStatus: Failed" \
            --quiet
    fi
}

# Weekly backup report
weekly_backup_report() {
    echo "Weekly Backup Report:"
    echo "===================="
    
    # Successful backups
    echo "Successful Backups:"
    trello-cli card list --list "$DAILY_BACKUP_LIST" --fields name --format json | \
        jq -r '.[] | .name'
    
    # Failed backups
    echo -e "\nFailed Backups:"
    trello-cli card list --list "$FAILED_BACKUP_LIST" --fields name --format json | \
        jq -r '.[] | .name'
}
```

## Customer Support

### Support Ticket Management

```bash
#!/bin/bash
# Customer support ticket automation

SUPPORT_BOARD="support-board-id"
NEW_TICKETS_LIST="new-tickets-id"
IN_PROGRESS_LIST="in-progress-id"
RESOLVED_LIST="resolved-id"

# Process incoming support tickets
process_support_ticket() {
    local ticket_id="$1"
    local customer_name="$2"
    local issue_description="$3"
    local priority="$4"
    
    # Determine list based on priority
    case "$priority" in
        "high"|"critical")
            TARGET_LIST="$NEW_TICKETS_LIST"
            ;;
        "medium"|"low")
            TARGET_LIST="$NEW_TICKETS_LIST"
            ;;
    esac
    
    # Create support card
    trello-cli card create --list "$TARGET_LIST" \
        "Ticket #$ticket_id - $customer_name" \
        --desc "Customer: $customer_name\nIssue: $issue_description\nPriority: $priority\nStatus: New" \
        --quiet
}

# Support metrics
generate_support_metrics() {
    echo "Support Metrics:"
    echo "================"
    
    # New tickets
    NEW_COUNT=$(trello-cli card list --list "$NEW_TICKETS_LIST" --format json | jq 'length')
    echo "New Tickets: $NEW_COUNT"
    
    # In progress
    IN_PROGRESS_COUNT=$(trello-cli card list --list "$IN_PROGRESS_LIST" --format json | jq 'length')
    echo "In Progress: $IN_PROGRESS_COUNT"
    
    # Resolved (last 24 hours)
    RESOLVED_COUNT=$(trello-cli card list --list "$RESOLVED_LIST" --format json | jq 'length')
    echo "Resolved: $RESOLVED_COUNT"
}
```

## Event Management

### Conference Planning

```bash
#!/bin/bash
# Conference and event planning automation

EVENT_BOARD="conference-board-id"
PLANNING_LIST="planning-id"
LOGISTICS_LIST="logistics-id"
SPEAKERS_LIST="speakers-id"
SPONSORS_LIST="sponsors-id"

# Plan conference
plan_conference() {
    local event_name="$1"
    local event_date="$2"
    local venue="$3"
    
    echo "Planning: $event_name"
    
    # Planning phase tasks
    trello-cli card create --list "$PLANNING_LIST" \
        "Define conference theme" \
        --desc "Event: $event_name\nDue: 3 months before" \
        --quiet
    
    trello-cli card create --list "$PLANNING_LIST" \
        "Create event website" \
        --desc "Event: $event_name\nDue: 2 months before" \
        --quiet
    
    # Logistics tasks
    trello-cli card create --list "$LOGISTICS_LIST" \
        "Book venue: $venue" \
        --desc "Event: $event_name\nDate: $event_date\nDue: 6 months before" \
        --quiet
    
    trello-cli card create --list "$LOGISTICS_LIST" \
        "Arrange catering" \
        --desc "Event: $event_name\nDate: $event_date\nDue: 1 month before" \
        --quiet
}

# Speaker management
manage_speakers() {
    local speaker_name="$1"
    local topic="$2"
    local contact_email="$3"
    
    trello-cli card create --list "$SPEAKERS_LIST" \
        "Speaker: $speaker_name" \
        --desc "Topic: $topic\nContact: $contact_email\nStatus: Confirmed" \
        --quiet
}
```
