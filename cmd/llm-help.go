package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	llmHelpCmd := &cobra.Command{
		Use:   "llm-help",
		Short: "LLM usage guidelines and best practices",
		Long:  "Output comprehensive usage guidelines optimized for LLM consumption with best practices, workflow patterns, and anti-patterns to avoid.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(llmHelpText)
		},
	}

	rootCmd.AddCommand(llmHelpCmd)
}

const llmHelpText = `
═══════════════════════════════════════════════════════════════════════════════
  TRELLO-CLI: LLM USAGE GUIDELINES
═══════════════════════════════════════════════════════════════════════════════

⚡ QUICK START WORKFLOW
────────────────────────────────────────────────────────────────────────────────
1. Run 'trello-cli schema' ONCE to understand all available commands
2. Use field filtering (--fields) to retrieve only necessary data
3. Set max-tokens (--max-tokens) to stay within context limits
4. Prefer batch operations for multiple modifications
5. Use --format json for parsing, --format markdown for user display

📋 CORE PRINCIPLES
────────────────────────────────────────────────────────────────────────────────

1️⃣  SCHEMA FIRST
   Always start by fetching the schema to understand capabilities:

   $ trello-cli schema

   Cache this for the conversation. Don't re-fetch for every command.

2️⃣  FILTER FIELDS
   ALWAYS use --fields to request only what you need:

   ✓ GOOD:   trello-cli card list --list <id> --fields name,desc,due
   ✗ AVOID:  trello-cli card list --list <id> --verbose

   Common field combinations:
   • Overview:     --fields name,desc,closed
   • Tasks:        --fields name,desc,due,labels
   • Minimal:      --fields name

3️⃣  SET TOKEN LIMITS
   Use --max-tokens to prevent context overflow:

   $ trello-cli board get <id> --max-tokens 2000

   Recommended limits by use case:
   • Quick overview:  500-1000 tokens
   • Task list:       1500-2000 tokens
   • Detailed view:   3000-4000 tokens

4️⃣  BATCH OPERATIONS
   When creating/modifying multiple items, use batch:

   ✓ GOOD:   echo '{"operations":[...]}' | trello-cli batch stdin
   ✗ AVOID:  Multiple sequential trello-cli card create commands

   Batch format:
   {
     "operations": [
       {
         "type": "card",
         "resource": "card",
         "action": "create",
         "data": {"name": "Task 1", "list_id": "abc123"}
       }
     ],
     "continue_on_error": true
   }

5️⃣  CHOOSE OUTPUT FORMAT
   • --format json:     For parsing and processing
   • --format markdown: For user-facing summaries
   • --quiet:           For automation (returns only IDs)

📖 RECOMMENDED WORKFLOW PATTERN
────────────────────────────────────────────────────────────────────────────────
# Step 1: Discover (once per conversation)
$ trello-cli schema

# Step 2: Query with optimization
$ trello-cli board list --fields name,desc,closed --max-tokens 1000

# Step 3: Get specific data
$ trello-cli card list --list <id> --fields name,desc,due --max-tokens 2000

# Step 4: Batch modifications
$ echo '{"operations":[...]}' | trello-cli batch stdin --quiet

# Step 5: Present results to user
$ trello-cli board get <id> --format markdown

🎯 CONTEXT OPTIMIZATION LEVELS
────────────────────────────────────────────────────────────────────────────────
Choose the right detail level for your task:

OVERVIEW (500 tokens)
$ trello-cli board list --fields name,desc,closed --max-tokens 500

TASK LIST (1500 tokens)
$ trello-cli card list --list <id> --fields name,desc,due,labels --max-tokens 1500

FULL DETAILS (use sparingly)
$ trello-cli card get <id> --verbose

⚠️  ANTI-PATTERNS TO AVOID
────────────────────────────────────────────────────────────────────────────────
❌ Fetching verbose data when specific fields would suffice
❌ Making sequential API calls instead of using batch operations
❌ Not setting --max-tokens limits (risks context overflow)
❌ Re-fetching schema for every command
❌ Using --verbose by default
❌ Requesting all boards/cards without field filtering
❌ Not using --quiet in multi-step automation

✅ BEST PRACTICES CHECKLIST
────────────────────────────────────────────────────────────────────────────────
☐ Fetch schema once at conversation start
☐ Always specify --fields with minimal required fields
☐ Set --max-tokens appropriate to context budget
☐ Use batch operations for multiple creates/updates
☐ Use --format json when processing data
☐ Use --format markdown when showing results to users
☐ Use --quiet when you only need IDs for follow-up operations
☐ Enable --debug only when troubleshooting
☐ Include continue_on_error: true in batch operations

🔧 COMMON OPERATIONS
────────────────────────────────────────────────────────────────────────────────
# List boards with minimal data
trello-cli board list --fields name,id --max-tokens 800

# Get board structure
trello-cli board get <board-id> --fields name,desc,lists --max-tokens 1500

# List cards efficiently
trello-cli card list --list <list-id> --fields name,desc,due --max-tokens 1200

# Create multiple cards
echo '{
  "operations": [
    {"type":"card","resource":"card","action":"create","data":{"name":"Task 1","list_id":"<id>"}},
    {"type":"card","resource":"card","action":"create","data":{"name":"Task 2","list_id":"<id>"}}
  ],
  "continue_on_error": true
}' | trello-cli batch stdin --format json

# Get card with checklists
trello-cli card get <card-id> --fields name,desc,checklists --max-tokens 1000

📚 ADDITIONAL RESOURCES
────────────────────────────────────────────────────────────────────────────────
• Schema:        trello-cli schema
• Main help:     trello-cli --help
• Command help:  trello-cli <command> --help
• Batch docs:    trello-cli batch --help
• Config:        trello-cli config show

═══════════════════════════════════════════════════════════════════════════════
💡 TIP: Bookmark this workflow: schema → filter → batch → present
═══════════════════════════════════════════════════════════════════════════════
`
