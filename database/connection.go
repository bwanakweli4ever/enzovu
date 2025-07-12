package database

import (
	"database/sql"
	"fmt"
	"time"

	"enzovu/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"           // PostgreSQL
	_ "github.com/mattn/go-sqlite3" // SQLite
)

var DB *sql.DB

type ConnectionConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Connect() error {
	cfg := config.GetConfig()

	var dsn string
	var err error

	switch cfg.Database.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
			cfg.Database.Charset,
		)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Database,
		)
	case "sqlite3":
		dsn = cfg.Database.Database + ".db"
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	DB, err = sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	connConfig := getConnectionConfig()
	DB.SetMaxOpenConns(connConfig.MaxOpenConns)
	DB.SetMaxIdleConns(connConfig.MaxIdleConns)
	DB.SetConnMaxLifetime(connConfig.ConnMaxLifetime)
	DB.SetConnMaxIdleTime(connConfig.ConnMaxIdleTime)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Printf("✅ Database connected successfully (%s)\n", cfg.Database.Driver)
	return nil
}

func getConnectionConfig() ConnectionConfig {
	cfg := config.GetConfig()

	if cfg.App.Environment == "production" {
		return ConnectionConfig{
			MaxOpenConns:    25,
			MaxIdleConns:    25,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
		}
	}

	return ConnectionConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    10,
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}
}

func Close() error {
	if DB != nil {
		fmt.Println("🔌 Closing database connection...")
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}

// Health check for the database connection
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}

// Transaction helper.
func Transaction(fn func(*sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
