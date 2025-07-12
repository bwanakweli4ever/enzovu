package database

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
)

type Migration interface {
	Up(*sql.DB) error
	Down(*sql.DB) error
	GetName() string
	GetTimestamp() string
}

type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		timestamp VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getMigratedNames() (map[string]bool, error) {
	migrated := make(map[string]bool)

	rows, err := m.db.Query("SELECT name FROM migrations")
	if err != nil {
		return migrated, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return migrated, err
		}
		migrated[name] = true
	}

	return migrated, nil
}

func (m *Migrator) recordMigration(migration Migration) error {
	_, err := m.db.Exec(
		"INSERT INTO migrations (name, timestamp) VALUES (?, ?)",
		migration.GetName(),
		migration.GetTimestamp(),
	)
	return err
}

func (m *Migrator) removeMigrationRecord(migration Migration) error {
	_, err := m.db.Exec("DELETE FROM migrations WHERE name = ?", migration.GetName())
	return err
}

func (m *Migrator) Migrate() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrated, err := m.getMigratedNames()
	if err != nil {
		return fmt.Errorf("failed to get migrated names: %w", err)
	}

	// Sort migrations by timestamp
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].GetTimestamp() < m.migrations[j].GetTimestamp()
	})

	migratedCount := 0
	for _, migration := range m.migrations {
		if migrated[migration.GetName()] {
			continue // Already migrated
		}

		log.Printf("⬆️  Running migration: %s", migration.GetName())

		if err := migration.Up(m.db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.GetName(), err)
		}

		if err := m.recordMigration(migration); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.GetName(), err)
		}

		migratedCount++
	}

	if migratedCount > 0 {
		fmt.Printf("✅ Ran %d migrations successfully\n", migratedCount)
	} else {
		fmt.Println("ℹ️  No new migrations to run")
	}

	return nil
}

func (m *Migrator) Rollback(steps int) error {
	if steps <= 0 {
		steps = 1
	}

	migrated, err := m.getMigratedNames()
	if err != nil {
		return fmt.Errorf("failed to get migrated names: %w", err)
	}

	// Get migrations to rollback (reverse order)
	var toRollback []Migration
	for i := len(m.migrations) - 1; i >= 0 && len(toRollback) < steps; i-- {
		migration := m.migrations[i]
		if migrated[migration.GetName()] {
			toRollback = append(toRollback, migration)
		}
	}

	for _, migration := range toRollback {
		log.Printf("⬇️  Rolling back migration: %s", migration.GetName())

		if err := migration.Down(m.db); err != nil {
			return fmt.Errorf("rollback of %s failed: %w", migration.GetName(), err)
		}

		if err := m.removeMigrationRecord(migration); err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", migration.GetName(), err)
		}
	}

	fmt.Printf("✅ Rolled back %d migrations successfully\n", len(toRollback))
	return nil
}

func (m *Migrator) Status() error {
	migrated, err := m.getMigratedNames()
	if err != nil {
		return fmt.Errorf("failed to get migrated names: %w", err)
	}

	fmt.Println("\n📊 Migration Status:")
	fmt.Println(strings.Repeat("-", 50))

	for _, migration := range m.migrations {
		status := "❌ Not migrated"
		if migrated[migration.GetName()] {
			status = "✅ Migrated"
		}
		fmt.Printf("%-30s %s\n", migration.GetName(), status)
	}

	fmt.Printf("\nTotal: %d migrations, %d migrated\n", len(m.migrations), len(migrated))
	return nil
}
