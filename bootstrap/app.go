package bootstrap

import (
	"enzovu/config"
	"fmt"
)

func InitializeApp() {
	fmt.Println("Initializing Enzovu Framework...")
	config.LoadConfig() // Load configurations
}
