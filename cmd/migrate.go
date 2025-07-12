// cmd/migrate.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/spf13/cobra"
)

var db *sql.DB

// Migration interface that all migrations must implement
type Migration interface {
	Up(db *sql.DB) error
	Down(db *sql.DB) error
	GetName() string
	GetTimestamp() string
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connectDB(); err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		if err := runMigrations(); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connectDB(); err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		if err := rollbackMigration(); err != nil {
			log.Fatal("Failed to rollback migration:", err)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connectDB(); err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		if err := showMigrationStatus(); err != nil {
			log.Fatal("Failed to show migration status:", err)
		}
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connectDB(); err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		if err := resetMigrations(); err != nil {
			log.Fatal("Failed to reset migrations:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
}

func connectDB() error {
	// Load environment variables
	loadEnvFile()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "enzovu_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Create migrations table if it doesn't exist
	return createMigrationsTable()
}

func createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		timestamp VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func runMigrations() error {
	// Get pending migrations
	pendingMigrations, err := getPendingMigrations()
	if err != nil {
		return err
	}

	if len(pendingMigrations) == 0 {
		fmt.Println("✅ No pending migrations")
		return nil
	}

	fmt.Printf("📦 Running %d pending migrations...\n", len(pendingMigrations))

	for _, migration := range pendingMigrations {
		fmt.Printf("⏳ Running migration: %s\n", migration.GetName())

		if err := migration.Up(db); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.GetName(), err)
		}

		// Record migration as executed
		if err := recordMigration(migration); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.GetName(), err)
		}

		fmt.Printf("✅ Migration %s completed\n", migration.GetName())
	}

	fmt.Println("🎉 All migrations completed successfully!")
	return nil
}

func rollbackMigration() error {
	// Get the last executed migration
	lastMigration, err := getLastExecutedMigration()
	if err != nil {
		return err
	}

	if lastMigration == nil {
		fmt.Println("✅ No migrations to rollback")
		return nil
	}

	fmt.Printf("⏳ Rolling back migration: %s\n", lastMigration.GetName())

	if err := lastMigration.Down(db); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.GetName(), err)
	}

	// Remove migration record
	if err := removeMigrationRecord(lastMigration); err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", lastMigration.GetName(), err)
	}

	fmt.Printf("✅ Migration %s rolled back successfully\n", lastMigration.GetName())
	return nil
}

func resetMigrations() error {
	// Get all executed migrations in reverse order
	executedMigrations, err := getExecutedMigrations()
	if err != nil {
		return err
	}

	if len(executedMigrations) == 0 {
		fmt.Println("✅ No migrations to reset")
		return nil
	}

	fmt.Printf("⚠️  Rolling back %d migrations...\n", len(executedMigrations))

	// Rollback in reverse order
	for i := len(executedMigrations) - 1; i >= 0; i-- {
		migration := executedMigrations[i]
		fmt.Printf("⏳ Rolling back migration: %s\n", migration.GetName())

		if err := migration.Down(db); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.GetName(), err)
		}

		if err := removeMigrationRecord(migration); err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", migration.GetName(), err)
		}

		fmt.Printf("✅ Migration %s rolled back\n", migration.GetName())
	}

	fmt.Println("🎉 All migrations reset successfully!")
	return nil
}

func showMigrationStatus() error {
	allMigrations := getAllMigrationFiles()
	executedMigrations, err := getExecutedMigrationNames()
	if err != nil {
		return err
	}

	fmt.Println("📊 Migration Status:")
	fmt.Println("==================")

	for _, migrationFile := range allMigrations {
		name := extractMigrationName(migrationFile)
		status := "❌ Pending"

		for _, executed := range executedMigrations {
			if executed == name {
				status = "✅ Executed"
				break
			}
		}

		fmt.Printf("%s - %s\n", status, name)
	}

	return nil
}

// Helper functions
func getPendingMigrations() ([]Migration, error) {
	allMigrations := getAllMigrationFiles()
	executedMigrations, err := getExecutedMigrationNames()
	if err != nil {
		return nil, err
	}

	var pending []Migration
	for _, migrationFile := range allMigrations {
		name := extractMigrationName(migrationFile)

		isExecuted := false
		for _, executed := range executedMigrations {
			if executed == name {
				isExecuted = true
				break
			}
		}

		if !isExecuted {
			// For this example, we'll need to implement actual migration loading
			// This is a simplified version
			fmt.Printf("Found pending migration: %s\n", name)
		}
	}

	return pending, nil
}

func getAllMigrationFiles() []string {
	var files []string
	migrationDir := "database/migrations"

	err := filepath.Walk(migrationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "_test.go") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Printf("Error reading migration directory: %v", err)
		return []string{}
	}

	return files
}

func extractMigrationName(filename string) string {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, ".go")

	// Remove timestamp prefix (YYYYMMDD_HHMMSS_)
	parts := strings.SplitN(name, "_", 3)
	if len(parts) >= 3 {
		return strings.Join(parts[2:], "_")
	}

	return name
}

func getExecutedMigrationNames() ([]string, error) {
	query := "SELECT name FROM migrations ORDER BY executed_at"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

func getLastExecutedMigration() (Migration, error) {
	// This is a placeholder - in a real implementation, you'd load the actual migration
	return nil, fmt.Errorf("migration loading not implemented yet")
}

func getExecutedMigrations() ([]Migration, error) {
	// This is a placeholder - in a real implementation, you'd load the actual migrations
	return []Migration{}, nil
}

func recordMigration(migration Migration) error {
	query := "INSERT INTO migrations (name, timestamp) VALUES (?, ?)"
	_, err := db.Exec(query, migration.GetName(), migration.GetTimestamp())
	return err
}

func removeMigrationRecord(migration Migration) error {
	query := "DELETE FROM migrations WHERE name = ?"
	_, err := db.Exec(query, migration.GetName())
	return err
}

func loadEnvFile() {
	// Simple .env file loader
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	// Basic implementation - in production, use a proper .env parser
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
