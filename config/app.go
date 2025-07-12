package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Environment string
	Port        string
	Debug       bool
	Name        string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

var AppConf *Config

func LoadConfig() *Config {
	loadEnvFile()

	config := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8000"),
			Debug:       getEnvBool("APP_DEBUG", true),
			Name:        getEnv("APP_NAME", "Enzovu App"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "mysql"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "enzovu_db"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
	}

	AppConf = config
	fmt.Printf("🔧 Configuration loaded - Environment: %s, Port: %s\n",
		config.App.Environment, config.App.Port)

	return config
}

func loadEnvFile() {
	envFiles := []string{".env.local", ".env"}

	for _, filename := range envFiles {
		file, err := os.Open(filename)
		if err != nil {
			continue // Try next file
		}
		defer file.Close()

		fmt.Printf("📄 Loading environment from %s\n", filename)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Remove quotes if present
				if len(value) >= 2 &&
					((value[0] == '"' && value[len(value)-1] == '"') ||
						(value[0] == '\'' && value[len(value)-1] == '\'')) {
					value = value[1 : len(value)-1]
				}

				os.Setenv(key, value)
			}
		}
		break // Successfully loaded one file
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// GetConfig returns the global configuration.
func GetConfig() *Config {
	if AppConf == nil {
		return LoadConfig()
	}
	return AppConf
}

// IsProduction checks if the app is running in production.
func IsProduction() bool {
	return GetConfig().App.Environment == "production"
}

// IsDevelopment checks if the app is running in development
func IsDevelopment() bool {
	return GetConfig().App.Environment == "development"
}

// IsDebug checks if debug mode is enabled
func IsDebug() bool {
	return GetConfig().App.Debug
}
