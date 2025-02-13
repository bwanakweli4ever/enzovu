package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// CreateCmd: Main command for generating resources (models, controllers, etc.)
var CreateCmd = &cobra.Command{
	Use:   "create [resource] [name]",
	Short: "Create a new resource (model, controller, etc.)",
	Long:  `Create a new resource like a model, controller, or migration for the Enzovu framework.`,
	Args:  cobra.ExactArgs(2), // Ensure exactly two arguments: resource type and name
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		resourceName := args[1]

		switch resourceType {
		case "model":
			createModel(resourceName)
		case "controller":
			createController(resourceName)
		default:
			fmt.Println("Invalid resource type. Use 'model' or 'controller'.")
		}
	},
}

// Function to create a new model file
func createModel(name string) {
	// Define the model path
	modelPath := fmt.Sprintf("app/Models/%s.go", strings.ToLower(name))
	file, err := os.Create(modelPath)
	if err != nil {
		fmt.Println("Error creating model:", err)
		return
	}
	defer file.Close()

	// Generic model content template
	modelContent := fmt.Sprintf(`package models

// %[1]s represents the model for the %[1]s resource.
type %[1]s struct {
	// ID is the primary identifier for the resource
	ID    int    `+"`json:\"id\"`"+`

	// Name is the name of the resource
	Name  string `+"`json:\"name\"`"+`

	// Email is the email of the resource owner (optional)
	Email string `+"`json:\"email\"`"+`

	// Add more fields as necessary (e.g., Description, CreatedAt, UpdatedAt)
}

// Get%[1]s returns a sample %[1]s instance
func Get%[1]s() *%[1]s {
	return &%[1]s{
		ID:    1,
		Name:  "Example %[1]s",
		Email: "example@email.com",
	}
}
`, name)

	// Write the generated model content to the file
	_, err = file.WriteString(modelContent)
	if err != nil {
		fmt.Println("Error writing to model file:", err)
		return
	}

	// Output success message
	fmt.Printf("✅ Model %s created successfully at %s\n", name, modelPath)
}

// Function to create a new controller file
func createController(name string) {
	controllerPath := fmt.Sprintf("app/Http/Controllers/%s_controller.go", strings.ToLower(name))
	file, err := os.Create(controllerPath)
	if err != nil {
		fmt.Println("Error creating controller:", err)
		return
	}
	defer file.Close()

	// Define basic controller structure
	controllerContent := fmt.Sprintf(`package controllers

import "net/http"

// %s handles HTTP requests for %s
func %s(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from %s controller!"))
}
`, name, name, name, name)

	_, err = file.WriteString(controllerContent)
	if err != nil {
		fmt.Println("Error writing to controller file:", err)
		return
	}

	fmt.Printf("✅ Controller %s created successfully at %s\n", name, controllerPath)
}
