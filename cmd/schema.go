package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// CommandSchema represents the schema for a command
type CommandSchema struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Usage       string                    `json:"usage"`
	Subcommands []SubcommandSchema        `json:"subcommands,omitempty"`
	GlobalFlags []FlagSchema              `json:"global_flags,omitempty"`
	Examples    []string                  `json:"examples,omitempty"`
}

// SubcommandSchema represents a subcommand
type SubcommandSchema struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Usage       string       `json:"usage"`
	Arguments   []ArgSchema  `json:"arguments,omitempty"`
	Flags       []FlagSchema `json:"flags,omitempty"`
	Examples    []string     `json:"examples,omitempty"`
}

// ArgSchema represents a command argument
type ArgSchema struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
}

// FlagSchema represents a command flag
type FlagSchema struct {
	Name        string `json:"name"`
	Short       string `json:"short,omitempty"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required"`
}

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Output complete CLI schema in JSON format",
		Long:  "Output a comprehensive JSON schema of all commands, subcommands, flags, and arguments for LLM consumption.",
		RunE: func(cmd *cobra.Command, args []string) error {
			schema := buildSchema()

			output, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal schema: %w", err)
			}

			fmt.Println(string(output))
			return nil
		},
	}

	rootCmd.AddCommand(schemaCmd)
}

func buildSchema() CommandSchema {
	return CommandSchema{
		Name:        "trello-cli",
		Description: "A comprehensive Trello CLI tool built in Go that provides full access to Trello's API with features optimized for LLM integration",
		Usage:       "trello-cli [command]",
		GlobalFlags: []FlagSchema{
			{Name: "api-key", Description: "Trello API key (overrides env/config)", Type: "string", Required: false},
			{Name: "token", Description: "Trello token (overrides env/config)", Type: "string", Required: false},
			{Name: "format", Short: "f", Description: "Output format (json, markdown)", Type: "string", Default: "markdown", Required: false},
			{Name: "fields", Description: "Specific fields to include in output", Type: "[]string", Required: false},
			{Name: "max-tokens", Description: "Maximum tokens in output (0 = unlimited)", Type: "int", Default: "0", Required: false},
			{Name: "verbose", Short: "v", Description: "Verbose output", Type: "bool", Default: "false", Required: false},
			{Name: "quiet", Short: "q", Description: "Quiet mode (minimal output)", Type: "bool", Default: "false", Required: false},
			{Name: "debug", Description: "Debug mode (show API calls)", Type: "bool", Default: "false", Required: false},
		},
		Subcommands: []SubcommandSchema{
			// Board commands
			{
				Name:        "board list",
				Description: "List all boards accessible to the authenticated user",
				Usage:       "trello-cli board list [flags]",
				Examples:    []string{"trello-cli board list", "trello-cli board list --format json", "trello-cli board list --fields name,desc,url"},
			},
			{
				Name:        "board get",
				Description: "Get detailed information about a specific board",
				Usage:       "trello-cli board get <board-id> [flags]",
				Arguments:   []ArgSchema{{Name: "board-id", Description: "ID of the board to retrieve", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli board get 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli board get <board-id> --format json"},
			},
			{
				Name:        "board create",
				Description: "Create a new Trello board",
				Usage:       "trello-cli board create <name> [flags]",
				Arguments:   []ArgSchema{{Name: "name", Description: "Name of the board to create", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "desc", Description: "Board description", Type: "string", Required: false}},
				Examples:    []string{"trello-cli board create \"My New Board\"", "trello-cli board create \"Project Board\" --desc \"Board for project management\""},
			},
			{
				Name:        "board delete",
				Description: "Delete a board permanently",
				Usage:       "trello-cli board delete <board-id> [flags]",
				Arguments:   []ArgSchema{{Name: "board-id", Description: "ID of the board to delete", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli board delete 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},
			{
				Name:        "board add-member",
				Description: "Add a member to a board",
				Usage:       "trello-cli board add-member <board-id> <email> [flags]",
				Arguments: []ArgSchema{
					{Name: "board-id", Description: "ID of the board", Required: true, Type: "string"},
					{Name: "email", Description: "Email address of the member to add", Required: true, Type: "string"},
				},
				Examples: []string{"trello-cli board add-member 5f8b8c8d8e8f8a8b8c8d8e8f user@example.com"},
			},

			// Card commands
			{
				Name:        "card list",
				Description: "List all cards in a specific list",
				Usage:       "trello-cli card list --list <list-id> [flags]",
				Flags:       []FlagSchema{{Name: "list", Description: "ID of the list to list cards from", Type: "string", Required: true}},
				Examples:    []string{"trello-cli card list --list 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli card list --list <list-id> --fields name,desc,due"},
			},
			{
				Name:        "card get",
				Description: "Get detailed information about a specific card",
				Usage:       "trello-cli card get <card-id> [flags]",
				Arguments:   []ArgSchema{{Name: "card-id", Description: "ID of the card to retrieve", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli card get 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli card get <card-id> --format json"},
			},
			{
				Name:        "card create",
				Description: "Create a new card in a specific list",
				Usage:       "trello-cli card create --list <list-id> <name> [flags]",
				Arguments:   []ArgSchema{{Name: "name", Description: "Name of the card to create", Required: true, Type: "string"}},
				Flags: []FlagSchema{
					{Name: "list", Description: "ID of the list to create the card in", Type: "string", Required: true},
					{Name: "desc", Description: "Card description", Type: "string", Required: false},
				},
				Examples: []string{"trello-cli card create --list 5f8b8c8d8e8f8a8b8c8d8e8f \"My New Card\"", "trello-cli card create --list <list-id> \"Task Card\" --desc \"Description of the task\""},
			},
			{
				Name:        "card move",
				Description: "Move a card to another list",
				Usage:       "trello-cli card move <card-id> --list <target-list-id> [flags]",
				Arguments:   []ArgSchema{{Name: "card-id", Description: "ID of the card to move", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "list", Description: "ID of the target list", Type: "string", Required: true}},
				Examples:    []string{"trello-cli card move 5f8b8c8d8e8f8a8b8c8d8e8f --list 5f8b8c8d8e8f8a8b8c8d8e8g"},
			},
			{
				Name:        "card copy",
				Description: "Copy a card to another list",
				Usage:       "trello-cli card copy <card-id> --list <target-list-id> [flags]",
				Arguments:   []ArgSchema{{Name: "card-id", Description: "ID of the card to copy", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "list", Description: "ID of the target list", Type: "string", Required: true}},
				Examples:    []string{"trello-cli card copy 5f8b8c8d8e8f8a8b8c8d8e8f --list 5f8b8c8d8e8f8a8b8c8d8e8g"},
			},
			{
				Name:        "card archive",
				Description: "Archive a card (soft delete)",
				Usage:       "trello-cli card archive <card-id> [flags]",
				Arguments:   []ArgSchema{{Name: "card-id", Description: "ID of the card to archive", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli card archive 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},
			{
				Name:        "card delete",
				Description: "Permanently delete a card",
				Usage:       "trello-cli card delete <card-id> [flags]",
				Arguments:   []ArgSchema{{Name: "card-id", Description: "ID of the card to delete", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli card delete 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},

			// List commands
			{
				Name:        "list list",
				Description: "List all lists on a board",
				Usage:       "trello-cli list list --board <board-id> [flags]",
				Flags:       []FlagSchema{{Name: "board", Description: "ID of the board to list lists from", Type: "string", Required: true}},
				Examples:    []string{"trello-cli list list --board 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli list list --board <board-id> --fields name,closed"},
			},
			{
				Name:        "list get",
				Description: "Get detailed information about a specific list",
				Usage:       "trello-cli list get <list-id> [flags]",
				Arguments:   []ArgSchema{{Name: "list-id", Description: "ID of the list to retrieve", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli list get 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},
			{
				Name:        "list create",
				Description: "Create a new list on a board",
				Usage:       "trello-cli list create --board <board-id> <name> [flags]",
				Arguments:   []ArgSchema{{Name: "name", Description: "Name of the list to create", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "board", Description: "ID of the board to create the list on", Type: "string", Required: true}},
				Examples:    []string{"trello-cli list create --board 5f8b8c8d8e8f8a8b8c8d8e8f \"New List\""},
			},
			{
				Name:        "list archive",
				Description: "Archive a list (soft delete)",
				Usage:       "trello-cli list archive <list-id> [flags]",
				Arguments:   []ArgSchema{{Name: "list-id", Description: "ID of the list to archive", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli list archive 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},

			// Label commands
			{
				Name:        "label list",
				Description: "List all labels on a board",
				Usage:       "trello-cli label list --board <board-id> [flags]",
				Flags:       []FlagSchema{{Name: "board", Description: "ID of the board to list labels from", Type: "string", Required: true}},
				Examples:    []string{"trello-cli label list --board 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli label list --board <board-id> --fields name,color"},
			},
			{
				Name:        "label create",
				Description: "Create a new label on a board",
				Usage:       "trello-cli label create --board <board-id> --name <name> --color <color> [flags]",
				Flags: []FlagSchema{
					{Name: "board", Description: "ID of the board to create the label on", Type: "string", Required: true},
					{Name: "name", Description: "Name of the label", Type: "string", Required: true},
					{Name: "color", Description: "Color of the label (red, yellow, orange, green, blue, purple, pink, lime, sky, grey)", Type: "string", Required: true},
				},
				Examples: []string{"trello-cli label create --board 5f8b8c8d8e8f8a8b8c8d8e8f --name \"Important\" --color \"red\""},
			},
			{
				Name:        "label add",
				Description: "Add a label to a card",
				Usage:       "trello-cli label add <card-id> <label-id> [flags]",
				Arguments: []ArgSchema{
					{Name: "card-id", Description: "ID of the card to add the label to", Required: true, Type: "string"},
					{Name: "label-id", Description: "ID of the label to add", Required: true, Type: "string"},
				},
				Examples: []string{"trello-cli label add 5f8b8c8d8e8f8a8b8c8d8e8f 5f8b8c8d8e8f8a8b8c8d8e8g"},
			},

			// Checklist commands
			{
				Name:        "checklist list",
				Description: "List all checklists on a card",
				Usage:       "trello-cli checklist list --card <card-id> [flags]",
				Flags:       []FlagSchema{{Name: "card", Description: "ID of the card to list checklists from", Type: "string", Required: true}},
				Examples:    []string{"trello-cli checklist list --card 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli checklist list --card <card-id> --fields name,checkItems"},
			},
			{
				Name:        "checklist create",
				Description: "Create a new checklist on a card",
				Usage:       "trello-cli checklist create --card <card-id> <name> [flags]",
				Arguments:   []ArgSchema{{Name: "name", Description: "Name of the checklist to create", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "card", Description: "ID of the card to create the checklist on", Type: "string", Required: true}},
				Examples:    []string{"trello-cli checklist create --card 5f8b8c8d8e8f8a8b8c8d8e8f \"Task List\""},
			},
			{
				Name:        "checklist add-item",
				Description: "Add an item to a checklist",
				Usage:       "trello-cli checklist add-item <checklist-id> <item-name> [flags]",
				Arguments: []ArgSchema{
					{Name: "checklist-id", Description: "ID of the checklist to add the item to", Required: true, Type: "string"},
					{Name: "item-name", Description: "Name of the item to add", Required: true, Type: "string"},
				},
				Examples: []string{"trello-cli checklist add-item 5f8b8c8d8e8f8a8b8c8d8e8f \"Review code\""},
			},

			// Member commands
			{
				Name:        "member get",
				Description: "Get detailed information about a member",
				Usage:       "trello-cli member get <username-or-id> [flags]",
				Arguments:   []ArgSchema{{Name: "username-or-id", Description: "Username or ID of the member to retrieve", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli member get john_doe", "trello-cli member get me", "trello-cli member get 5f8b8c8d8e8f8a8b8c8d8e8f"},
			},
			{
				Name:        "member boards",
				Description: "List all boards that a member has access to",
				Usage:       "trello-cli member boards <username-or-id> [flags]",
				Arguments:   []ArgSchema{{Name: "username-or-id", Description: "Username or ID of the member", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli member boards john_doe", "trello-cli member boards me"},
			},

			// Attachment commands
			{
				Name:        "attachment list",
				Description: "List all attachments on a card",
				Usage:       "trello-cli attachment list --card <card-id> [flags]",
				Flags:       []FlagSchema{{Name: "card", Description: "ID of the card to list attachments from", Type: "string", Required: true}},
				Examples:    []string{"trello-cli attachment list --card 5f8b8c8d8e8f8a8b8c8d8e8f", "trello-cli attachment list --card <card-id> --fields name,url,mimeType"},
			},
			{
				Name:        "attachment add",
				Description: "Add an attachment to a card",
				Usage:       "trello-cli attachment add --card <card-id> <url> [flags]",
				Arguments:   []ArgSchema{{Name: "url", Description: "URL of the file to attach", Required: true, Type: "string"}},
				Flags:       []FlagSchema{{Name: "card", Description: "ID of the card to add the attachment to", Type: "string", Required: true}},
				Examples:    []string{"trello-cli attachment add --card 5f8b8c8d8e8f8a8b8c8d8e8f \"https://example.com/file.pdf\""},
			},

			// Batch commands
			{
				Name:        "batch file",
				Description: "Execute batch operations from a JSON or YAML file",
				Usage:       "trello-cli batch file <file-path> [flags]",
				Arguments:   []ArgSchema{{Name: "file-path", Description: "Path to the JSON file containing batch operations", Required: true, Type: "string"}},
				Examples:    []string{"trello-cli batch file operations.json", "trello-cli batch file operations.json --format json"},
			},
			{
				Name:        "batch stdin",
				Description: "Execute batch operations from JSON or YAML piped to stdin",
				Usage:       "trello-cli batch stdin [flags]",
				Examples:    []string{"cat operations.json | trello-cli batch stdin", "echo '{\"operations\":[...]}' | trello-cli batch stdin"},
			},

			// Config commands
			{
				Name:        "config show",
				Description: "Display the current configuration settings",
				Usage:       "trello-cli config show [flags]",
				Examples:    []string{"trello-cli config show"},
			},
			{
				Name:        "config set",
				Description: "Set configuration values for API credentials and default settings",
				Usage:       "trello-cli config set [flags]",
				Flags: []FlagSchema{
					{Name: "api-key", Description: "Trello API key", Type: "string", Required: false},
					{Name: "token", Description: "Trello token", Type: "string", Required: false},
					{Name: "default-format", Description: "Default output format", Type: "string", Default: "markdown", Required: false},
					{Name: "max-tokens", Description: "Default maximum tokens", Type: "int", Default: "4000", Required: false},
				},
				Examples: []string{"trello-cli config set --api-key \"key\" --token \"token\"", "trello-cli config set --default-format json"},
			},
			{
				Name:        "config path",
				Description: "Display the path to the configuration file",
				Usage:       "trello-cli config path [flags]",
				Examples:    []string{"trello-cli config path"},
			},
		},
		Examples: []string{
			"trello-cli board list",
			"trello-cli card create --list <list-id> \"New Card\" --desc \"Task description\"",
			"trello-cli board get <board-id> --format json --fields name,desc,url",
			"cat batch.json | trello-cli batch stdin",
		},
	}
}
