package commands

import (
	"github.com/spf13/cobra"
)

// Root command for go-craft
var rootCmd = &cobra.Command{
	Use:   "go-craft",
	Short: "Go-Craft is a CLI tool for the Enzovu framework",
	Long:  "Go-Craft helps you manage your Enzovu framework tasks like creating models, controllers, and more.",
}

// Initialize the CLI tool with subcommands
func init() {
	// Register subcommands for 'create', 'migrate', etc.

	// Example: rootCmd.AddCommand(createCmd)
}
