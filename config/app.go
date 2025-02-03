package config

import (
    "fmt"
    "os"
)

func LoadConfig() {
    appEnv := os.Getenv("APP_ENV")
    fmt.Println("App Environment: ", appEnv)
    // Further configuration loading here (e.g., database connection, etc.)
}
