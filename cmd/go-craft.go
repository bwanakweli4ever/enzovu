package main

import (
	"fmt"
	"os"

	"enzovu/app/commands" // Import the commands package

	"github.com/spf13/cobra"
)

// Root command for go-craft CLI
var rootCmd = &cobra.Command{
	Use:   "go-craft",
	Short: "Go-Craft is a CLI tool for the Enzovu framework",
	Long:  "Go-Craft helps you manage your Enzovu framework tasks like creating models, controllers, and more.",
}

// Initialize the CLI tool with subcommands
func init() {
	rootCmd.AddCommand(commands.CreateCmd) // Register the create command
}

// Main function to execute CLI commands
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
