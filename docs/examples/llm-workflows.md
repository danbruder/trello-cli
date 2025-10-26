# LLM Workflow Tutorials

Learn how to optimize Trello CLI for LLM integration with context optimization, batch processing, and intelligent output formatting.

## Context Optimization for LLMs

### Token-Aware Board Summaries

```bash
#!/bin/bash
# Get board context optimized for LLM token limits

BOARD_ID="your-board-id"
MAX_TOKENS=2000

# Get board summary with essential fields
trlo board get "$BOARD_ID" \
    --fields name,desc,url,closed \
    --format json \
    --max-tokens "$MAX_TOKENS"
```

### Intelligent Field Filtering

```bash
#!/bin/bash
# Get only essential fields for LLM processing

# For task management context
trlo card list --list "$LIST_ID" \
    --fields name,desc,due,labels \
    --format json \
    --max-tokens 1500

# For project overview context  
trlo board list \
    --fields name,desc,closed \
    --format json \
    --max-tokens 1000
```

### Context-Aware Data Extraction

```bash
#!/bin/bash
# Extract different levels of context based on use case

# High-level project overview
get_project_overview() {
    trlo board list --fields name,desc,closed --format json --max-tokens 800
}

# Detailed task context
get_task_context() {
    trlo card list --list "$1" --fields name,desc,labels,due --format json --max-tokens 1200
}

# Specific card details
get_card_details() {
    trlo card get "$1" --fields name,desc,labels,attachments --format json --max-tokens 500
}
```

## LLM-Generated Content Processing

### Process LLM Task Output

```bash
#!/bin/bash
# Process LLM-generated task list and create Trello cards

LLM_OUTPUT='{
  "tasks": [
    {"title": "Implement user authentication", "description": "Add login and registration", "priority": "high"},
    {"title": "Write API documentation", "description": "Document all endpoints", "priority": "medium"},
    {"title": "Add unit tests", "description": "Coverage for core functionality", "priority": "high"}
  ]
}'

# Convert LLM output to batch operations
echo "$LLM_OUTPUT" | jq '{
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
}' | trlo batch stdin --format json
```

### Dynamic List Selection Based on Priority

```bash
#!/bin/bash
# Route LLM-generated tasks to appropriate lists based on priority

process_llm_task() {
    local title="$1"
    local description="$2"
    local priority="$3"
    
    case "$priority" in
        "high"|"critical")
            LIST_ID="$HIGH_PRIORITY_LIST"
            ;;
        "medium")
            LIST_ID="$MEDIUM_PRIORITY_LIST"
            ;;
        "low")
            LIST_ID="$LOW_PRIORITY_LIST"
            ;;
        *)
            LIST_ID="$DEFAULT_LIST"
            ;;
    esac
    
    trlo card create --list "$LIST_ID" "$title" --desc "$description" --quiet
}

# Process LLM output
echo "$LLM_OUTPUT" | jq -r '.tasks[] | "\(.title)|\(.description)|\(.priority)"' | while IFS='|' read -r title desc priority; do
    process_llm_task "$title" "$desc" "$priority"
done
```

## Intelligent Workflow Automation

### LLM-Driven Project Setup

```bash
#!/bin/bash
# Use LLM to generate project structure and create Trello board

PROJECT_NAME="$1"
LLM_PROMPT="Create a project structure for: $PROJECT_NAME"

# Get LLM response (example using OpenAI API)
LLM_RESPONSE=$(curl -s -X POST "https://api.openai.com/v1/chat/completions" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "Content-Type: application/json" \
    -d "{
        \"model\": \"gpt-3.5-turbo\",
        \"messages\": [{\"role\": \"user\", \"content\": \"$LLM_PROMPT\"}]
    }" | jq -r '.choices[0].message.content')

# Create board
BOARD_ID=$(trlo board create "$PROJECT_NAME" --quiet)

# Parse LLM response and create lists
echo "$LLM_RESPONSE" | grep -E "^- " | sed 's/^- //' | while read -r list_name; do
    trlo list create --board "$BOARD_ID" "$list_name" --quiet
done
```

### Smart Task Categorization

```bash
#!/bin/bash
# Use LLM to categorize tasks and assign appropriate labels

categorize_task() {
    local task_title="$1"
    local task_description="$2"
    
    # Get LLM categorization
    CATEGORY=$(curl -s -X POST "https://api.openai.com/v1/chat/completions" \
        -H "Authorization: Bearer $OPENAI_API_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"model\": \"gpt-3.5-turbo\",
            \"messages\": [{\"role\": \"user\", \"content\": \"Categorize this task: $task_title - $task_description. Respond with only: bug, feature, documentation, or other\"}]
        }" | jq -r '.choices[0].message.content')
    
    # Create card and add appropriate label
    CARD_ID=$(trlo card create --list "$LIST_ID" "$task_title" --desc "$task_description" --quiet)
    
    # Get label ID based on category
    LABEL_ID=$(trlo label list --board "$BOARD_ID" --format json | jq -r ".[] | select(.name == \"$CATEGORY\") | .id")
    
    if [ "$LABEL_ID" != "null" ] && [ "$LABEL_ID" != "" ]; then
        trlo label add "$CARD_ID" "$LABEL_ID" --quiet
    fi
}
```

## Advanced LLM Integration Patterns

### Context-Aware Task Generation

```bash
#!/bin/bash
# Generate tasks based on current project context

get_project_context() {
    # Get current board state
    trlo board get "$BOARD_ID" --fields name,desc --format json
    trlo card list --list "$CURRENT_LIST" --fields name,desc,labels --format json --max-tokens 1000
}

generate_next_tasks() {
    local context="$1"
    
    # Use LLM to generate next tasks based on context
    curl -s -X POST "https://api.openai.com/v1/chat/completions" \
        -H "Authorization: Bearer $OPENAI_API_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"model\": \"gpt-3.5-turbo\",
            \"messages\": [{\"role\": \"user\", \"content\": \"Based on this project context: $context, suggest 3 next tasks. Return as JSON array with title and description fields.\"}]
        }" | jq -r '.choices[0].message.content'
}

# Get context and generate tasks
CONTEXT=$(get_project_context)
NEXT_TASKS=$(generate_next_tasks "$CONTEXT")

# Create tasks from LLM suggestions
echo "$NEXT_TASKS" | jq -r '.[] | "\(.title)|\(.description)"' | while IFS='|' read -r title desc; do
    trlo card create --list "$LIST_ID" "$title" --desc "$desc"
done
```

### Intelligent Progress Tracking

```bash
#!/bin/bash
# Use LLM to analyze project progress and suggest actions

analyze_progress() {
    # Get comprehensive project data
    BOARD_DATA=$(trlo board get "$BOARD_ID" --format json)
    CARDS_DATA=$(trlo card list --list "$LIST_ID" --fields name,desc,labels --format json --max-tokens 2000)
    
    # Combine data for LLM analysis
    ANALYSIS_DATA=$(echo "{\"board\": $BOARD_DATA, \"cards\": $CARDS_DATA}" | jq -c .)
    
    # Get LLM analysis
    curl -s -X POST "https://api.openai.com/v1/chat/completions" \
        -H "Authorization: Bearer $OPENAI_API_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"model\": \"gpt-3.5-turbo\",
            \"messages\": [{\"role\": \"user\", \"content\": \"Analyze this project data and suggest next actions: $ANALYSIS_DATA\"}]
        }" | jq -r '.choices[0].message.content'
}

# Run analysis and create action items
ANALYSIS=$(analyze_progress)
echo "Project Analysis: $ANALYSIS"
```

## Performance Optimization

### Efficient Data Fetching

```bash
#!/bin/bash
# Optimize data fetching for LLM processing

# Use field filtering to reduce token usage
get_optimized_context() {
    local context_type="$1"
    
    case "$context_type" in
        "overview")
            trlo board list --fields name,desc,closed --format json --max-tokens 500
            ;;
        "detailed")
            trlo card list --list "$LIST_ID" --fields name,desc,labels,due --format json --max-tokens 1500
            ;;
        "minimal")
            trlo card list --list "$LIST_ID" --fields name --format json --max-tokens 200
            ;;
    esac
}
```

### Batch Processing for LLM Output

```bash
#!/bin/bash
# Efficiently process large LLM outputs

process_llm_batch() {
    local llm_output="$1"
    
    # Convert to batch operations with error handling
    echo "$llm_output" | jq '{
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
    }' | trlo batch stdin --format json --quiet
}
```
