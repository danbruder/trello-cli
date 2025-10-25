package main

import "github.com/spf13/cobra"

// Global variables shared across all command files
// These are defined in main.go and used by all commands
var (
	format    string
	fields    []string
	maxTokens int
	verbose   bool
	quiet     bool
)

// Command variables
var (
	attachmentCmd *cobra.Command
	batchCmd      *cobra.Command
	boardCmd      *cobra.Command
	cardCmd       *cobra.Command
	checklistCmd  *cobra.Command
	configCmd     *cobra.Command
	labelCmd      *cobra.Command
	listCmd       *cobra.Command
	memberCmd     *cobra.Command
)

