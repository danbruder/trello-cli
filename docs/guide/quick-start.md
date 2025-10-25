# Quick Start

Let's get you up and running with Trello CLI in just a few steps.

## First Commands

### List Your Boards

```bash
trello-cli board list
```

This will show all boards you have access to in a formatted table.

### Get Board Details

```bash
trello-cli board get <board-id>
```

Replace `<board-id>` with an actual board ID from the previous command.

### List Cards in a Board

```bash
trello-cli card list --board <board-id>
```

## Basic Workflows

### Create a New Card

```bash
trello-cli card create --list <list-id> "My New Task"
```

### Move a Card

```bash
trello-cli card move <card-id> --list <target-list-id>
```

### Add a Label

```bash
trello-cli label add <card-id> <label-id>
```

## Output Formats

### Markdown (Default)

```bash
trello-cli board list
# Outputs formatted Markdown tables
```

### JSON

```bash
trello-cli board list --format json
# Outputs structured JSON
```

## LLM-Optimized Features

### Field Filtering

```bash
# Only include specific fields
trello-cli card list --list <list-id> --fields name,desc,due
```

### Token Limits

```bash
# Limit output to 2000 tokens
trello-cli board list --max-tokens 2000
```

### Quiet Mode for Scripting

```bash
# Minimal output for automation
trello-cli card create --list <list-id> "New Card" --quiet
```

## Next Steps

- Explore the [Command Reference](/reference/commands) for all available commands
- Check out [API Integration Examples](/examples/api-integration) for advanced usage
- Learn about [LLM Workflows](/examples/llm-workflows) for AI integration
