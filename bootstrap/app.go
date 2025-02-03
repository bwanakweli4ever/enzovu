package bootstrap

import (
	"enzovu/config"
	"fmt"
)

func InitializeApp() {
	fmt.Println("Initializing application...")
	config.LoadConfig() // Load configurations
}
