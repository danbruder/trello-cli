# Quick Start

Let's get you up and running with Trello CLI in just a few steps.

## Installation

If you haven't installed Trello CLI yet, the easiest way is via Homebrew:

```bash
brew tap danbruder/tap
brew install trlo
```

For other installation methods, see the [Installation Guide](/guide/installation).

## First Commands

### List Your Boards

```bash
trlo board list
```

This will show all boards you have access to in a formatted table.

### Get Board Details

```bash
trlo board get <board-id>
```

Replace `<board-id>` with an actual board ID from the previous command.

### List Cards in a List

```bash
trlo card list --list <list-id>
```

## Basic Workflows

### Create a New Card

```bash
trlo card create --list <list-id> "My New Task"
```

### Move a Card

```bash
trlo card move <card-id> --list <target-list-id>
```

### Add a Label

```bash
trlo label add <card-id> <label-id>
```

## Output Formats

### Markdown (Default)

```bash
trlo board list
# Outputs formatted Markdown tables
```

### JSON

```bash
trlo board list --format json
# Outputs structured JSON
```

## LLM-Optimized Features

### Field Filtering

```bash
# Only include specific fields
trlo card list --list <list-id> --fields name,desc,due
```

### Token Limits

```bash
# Limit output to 2000 tokens
trlo board list --max-tokens 2000
```

### Quiet Mode for Scripting

```bash
# Minimal output for automation
trlo card create --list <list-id> "New Card" --quiet
```

## Next Steps

- Explore the [Command Reference](/reference/commands) for all available commands
- Check out [API Integration Examples](/examples/api-integration) for advanced usage
- Learn about [LLM Workflows](/examples/llm-workflows) for AI integration
