// app/commands/create.go
package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// CreateCmd: Main command for generating resources (models, controllers, etc.)
var CreateCmd = &cobra.Command{
	Use:   "create [resource] [name]",
	Short: "Create a new resource (model, controller, middleware, migration)",
	Long:  `Create a new resource like a model, controller, middleware, or migration for the Enzovu framework.`,
	Args:  cobra.ExactArgs(2), // Ensure exactly two arguments: resource type and name
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		resourceName := args[1]

		switch resourceType {
		case "model":
			createModel(resourceName)
		case "controller":
			createController(resourceName)
		case "middleware":
			createMiddleware(resourceName)
		case "migration":
			createMigration(resourceName)
		default:
			fmt.Println("❌ Invalid resource type. Use 'model', 'controller', 'middleware', or 'migration'.")
		}
	},
}

// Function to create a new model file
func createModel(name string) {
	// Define the model path
	modelPath := fmt.Sprintf("app/Models/%s.go", strings.ToLower(name))

	// Check if file already exists
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		fmt.Printf("❌ Model %s already exists at %s\n", name, modelPath)
		return
	}

	file, err := os.Create(modelPath)
	if err != nil {
		fmt.Println("❌ Error creating model:", err)
		return
	}
	defer file.Close()

	// Fixed model content template
	modelContent := fmt.Sprintf(`package models

import (
	"encoding/json"
	"time"
)

// %[1]s represents the model for the %[1]s resource.
type %[1]s struct {
	// ID is the primary identifier for the resource
	ID    int    `+"`json:\"id\" db:\"id\"`"+`

	// Name is the name of the resource
	Name  string `+"`json:\"name\" db:\"name\"`"+`

	// Email is the email of the resource owner (optional)
	Email string `+"`json:\"email\" db:\"email\"`"+`

	// Timestamps
	CreatedAt time.Time `+"`json:\"created_at\" db:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\" db:\"updated_at\"`"+`

	// Add more fields as necessary
}

// Get%[1]s returns a sample %[1]s instance
func Get%[1]s() *%[1]s {
	return &%[1]s{
		ID:        1,
		Name:      "Example %[1]s",
		Email:     "example@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetAll%[1]s returns a slice of %[1]s instances
func GetAll%[1]s() []*%[1]s {
	return []*%[1]s{
		Get%[1]s(),
	}
}

// ToJSON converts the %[1]s to JSON
func (m *%[1]s) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// TableName returns the database table name for this model
func (%[1]s) TableName() string {
	return "%[2]s"
}

// Save saves the %[1]s to the database (placeholder)
func (m *%[1]s) Save() error {
	// TODO: Implement database save logic
	m.UpdatedAt = time.Now()
	return nil
}

// Delete removes the %[1]s from the database (placeholder)
func (m *%[1]s) Delete() error {
	// TODO: Implement database delete logic
	return nil
}
`, name, strings.ToLower(name)+"s")

	// Write the generated model content to the file
	_, err = file.WriteString(modelContent)
	if err != nil {
		fmt.Println("❌ Error writing to model file:", err)
		return
	}

	// Output success message
	fmt.Printf("✅ Model %s created successfully at %s\n", name, modelPath)
}

// Function to create a new controller file
func createController(name string) {
	controllerPath := fmt.Sprintf("app/Http/Controllers/%s_controller.go", strings.ToLower(name))

	// Check if file already exists
	if _, err := os.Stat(controllerPath); !os.IsNotExist(err) {
		fmt.Printf("❌ Controller %s already exists at %s\n", name, controllerPath)
		return
	}

	file, err := os.Create(controllerPath)
	if err != nil {
		fmt.Println("❌ Error creating controller:", err)
		return
	}
	defer file.Close()

	// Enhanced controller structure with CRUD operations
	controllerContent := fmt.Sprintf(`package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"enzovu/app/Models"
	"enzovu/routes"
)

// %[1]sController handles HTTP requests for %[1]s resources
type %[1]sController struct{}

// Index handles GET /%[2]s - List all %[2]s
func (c *%[1]sController) Index(w http.ResponseWriter, r *http.Request) {
	%[2]s := models.GetAll%[1]s()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(%[2]s)
}

// Show handles GET /%[2]s/{id} - Show a specific %[1]s
func (c *%[1]sController) Show(w http.ResponseWriter, r *http.Request) {
	// Get ID parameter from URL
	idStr := routes.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	// TODO: Fetch %[1]s by ID from database
	%[3]s := models.Get%[1]s()
	%[3]s.ID = id
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(%[3]s)
}

// Create handles POST /%[2]s - Create a new %[1]s
func (c *%[1]sController) Create(w http.ResponseWriter, r *http.Request) {
	var %[3]s models.%[1]s
	
	if err := json.NewDecoder(r.Body).Decode(&%[3]s); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// TODO: Validate and save to database
	if err := %[3]s.Save(); err != nil {
		http.Error(w, "Failed to create %[3]s", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(%[3]s)
}

// Update handles PUT /%[2]s/{id} - Update a %[1]s
func (c *%[1]sController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := routes.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var %[3]s models.%[1]s
	if err := json.NewDecoder(r.Body).Decode(&%[3]s); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	%[3]s.ID = id
	
	// TODO: Update in database
	if err := %[3]s.Save(); err != nil {
		http.Error(w, "Failed to update %[3]s", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(%[3]s)
}

// Delete handles DELETE /%[2]s/{id} - Delete a %[1]s
func (c *%[1]sController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := routes.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	// TODO: Fetch and delete from database
	%[3]s := models.Get%[1]s()
	%[3]s.ID = id
	
	if err := %[3]s.Delete(); err != nil {
		http.Error(w, "Failed to delete %[3]s", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Simple function-based handler for backward compatibility
func %[1]s(w http.ResponseWriter, r *http.Request) {
	controller := &%[1]sController{}
	controller.Index(w, r)
}
`, name, strings.ToLower(name)+"s", strings.ToLower(name))

	_, err = file.WriteString(controllerContent)
	if err != nil {
		fmt.Println("❌ Error writing to controller file:", err)
		return
	}

	fmt.Printf("✅ Controller %s created successfully at %s\n", name, controllerPath)
	fmt.Printf("💡 Don't forget to register your routes in routes/web.go\n")
}

// Function to create a new middleware file
func createMiddleware(name string) {
	middlewarePath := fmt.Sprintf("app/Http/Middleware/%sMiddleware.go", name)

	// Check if file already exists
	if _, err := os.Stat(middlewarePath); !os.IsNotExist(err) {
		fmt.Printf("❌ Middleware %s already exists at %s\n", name, middlewarePath)
		return
	}

	file, err := os.Create(middlewarePath)
	if err != nil {
		fmt.Println("❌ Error creating middleware:", err)
		return
	}
	defer file.Close()

	// Enhanced middleware template
	middlewareContent := fmt.Sprintf(`package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// %[1]sMiddleware provides %[2]s functionality
func %[1]sMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic
		start := time.Now()
		
		fmt.Printf("[%[1]s] Processing request: %%s %%s\n", r.Method, r.URL.Path)
		
		// TODO: Add your %[2]s logic here
		// Example checks:
		// - Authentication verification
		// - Rate limiting
		// - Request validation
		// - Logging
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Post-processing logic
		duration := time.Since(start)
		fmt.Printf("[%[1]s] Request completed in %%v\n", duration)
	})
}

// %[1]sConfig holds configuration for the %[1]s middleware
type %[1]sConfig struct {
	Enabled bool
	// Add configuration fields as needed
}

// Default%[1]sConfig returns default configuration
func Default%[1]sConfig() *%[1]sConfig {
	return &%[1]sConfig{
		Enabled: true,
	}
}

// %[1]sWithConfig creates middleware with custom configuration
func %[1]sWithConfig(config *%[1]sConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = Default%[1]sConfig()
	}
	
	return func(next http.Handler) http.Handler {
		if !config.Enabled {
			return next
		}
		return %[1]sMiddleware(next)
	}
}
`, name, strings.ToLower(name))

	_, err = file.WriteString(middlewareContent)
	if err != nil {
		fmt.Println("❌ Error writing to middleware file:", err)
		return
	}

	fmt.Printf("✅ Middleware %s created successfully at %s\n", name, middlewarePath)
	fmt.Printf("💡 Use it in your routes: router.Use(middleware.%sMiddleware)\n", name)
}

// Function to create a new migration file
func createMigration(name string) {
	timestamp := time.Now().Format("20060102_150405")
	migrationPath := fmt.Sprintf("database/migrations/%s_%s.go", timestamp, strings.ToLower(name))

	// Check if file already exists
	if _, err := os.Stat(migrationPath); !os.IsNotExist(err) {
		fmt.Printf("❌ Migration %s already exists at %s\n", name, migrationPath)
		return
	}

	file, err := os.Create(migrationPath)
	if err != nil {
		fmt.Println("❌ Error creating migration:", err)
		return
	}
	defer file.Close()

	// Enhanced migration template
	migrationContent := fmt.Sprintf(`package migrations

import (
	"database/sql"
	"fmt"
)

// %[1]s migration
type %[1]sMigration struct {
	Name      string
	Timestamp string
}

// NewMigration%[1]s creates a new migration instance
func NewMigration%[1]s() *%[1]sMigration {
	return &%[1]sMigration{
		Name:      "%[2]s",
		Timestamp: "%[3]s",
	}
}

// Up runs the migration
func (m *%[1]sMigration) Up(db *sql.DB) error {
	fmt.Printf("Running migration: %%s\n", m.Name)
	
	// TODO: Add your migration SQL here
	query := `+"`"+`
	-- Example: Create table
	-- CREATE TABLE example_table (
	--     id INT PRIMARY KEY AUTO_INCREMENT,
	--     name VARCHAR(255) NOT NULL,
	--     email VARCHAR(255) UNIQUE,
	--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	-- );
	`+"`"+`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to run migration %%s: %%w", m.Name, err)
	}
	
	fmt.Printf("✅ Migration %%s completed successfully\n", m.Name)
	return nil
}

// Down rolls back the migration
func (m *%[1]sMigration) Down(db *sql.DB) error {
	fmt.Printf("Rolling back migration: %%s\n", m.Name)
	
	// TODO: Add your rollback SQL here
	query := `+"`"+`
	-- Example: Drop table
	-- DROP TABLE IF EXISTS example_table;
	`+"`"+`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to rollback migration %%s: %%w", m.Name, err)
	}
	
	fmt.Printf("✅ Migration %%s rolled back successfully\n", m.Name)
	return nil
}

// GetName returns the migration name
func (m *%[1]sMigration) GetName() string {
	return m.Name
}

// GetTimestamp returns the migration timestamp
func (m *%[1]sMigration) GetTimestamp() string {
	return m.Timestamp
}
`, strings.Title(strings.ReplaceAll(name, "_", "")), name, timestamp)

	_, err = file.WriteString(migrationContent)
	if err != nil {
		fmt.Println("❌ Error writing to migration file:", err)
		return
	}

	fmt.Printf("✅ Migration %s created successfully at %s\n", name, migrationPath)
	fmt.Printf("💡 Register it in your migration runner to execute\n")
}
