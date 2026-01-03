package main

import (
	"os"

	"github.com/danbruder/trello-cli/cmd"
)

// Version information set during build (injected via ldflags)
var (
	version   = "1.1.0"
	buildTime = "unknown"
	goVersion = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.Version = version
	cmd.BuildTime = buildTime
	cmd.GoVersion = goVersion

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
