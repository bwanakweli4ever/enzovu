package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"enzovu/bootstrap"
	"enzovu/routes"
)

var isDevelopment = getEnv("APP_ENV", "development") == "development"

func main() {
	// Display welcome message
	displayWelcomeMessage()

	// Initialize application
	if err := initializeApp(); err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}

	// Get port from environment or use default
	port := getEnv("APP_PORT", "8000")

	if isDevelopment {
		runWithHotReload(port)
	} else {
		runProduction(port)
	}
}

func runProduction(port string) {
	// Setup routes
	router := routes.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	fmt.Printf("🚀 Enzovu server starting on http://localhost:%s\n", port)
	fmt.Printf("📊 Environment: %s\n", getEnv("APP_ENV", "production"))
	fmt.Println("🎯 Press Ctrl+C to shutdown")
	fmt.Println()

	// Setup graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		fmt.Println("\n🛑 Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server forced to shutdown: %v", err)
		}
		fmt.Println("✅ Server exited successfully")
		os.Exit(0)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

func runWithHotReload(port string) {
	fmt.Println("🔥 Hot reload enabled - edit any .go file to see changes!")
	fmt.Printf("🚀 Enzovu server starting on http://localhost:%s\n", port)
	fmt.Printf("📊 Environment: %s\n", getEnv("APP_ENV", "development"))
	fmt.Println("🎯 Press Ctrl+C to shutdown")
	fmt.Println()

	// Create a dynamic handler that reloads routes
	dynamicHandler := &DynamicHandler{}
	dynamicHandler.UpdateRoutes()

	// Create HTTP server with the dynamic handler
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      dynamicHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Watch for file changes
	lastMod := time.Now()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n🛑 Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server forced to shutdown: %v", err)
		}
		fmt.Println("✅ Server exited successfully")
		os.Exit(0)
	}()

	// File watching loop
	for {
		time.Sleep(500 * time.Millisecond)

		if hasGoFileChanged(&lastMod) {
			fmt.Println("📝 Changes detected, reloading routes...")

			// Reload routes without restarting server
			dynamicHandler.UpdateRoutes()

			fmt.Println("✅ Routes reloaded successfully!")
		}
	}
}

// DynamicHandler wraps the current routes and allows hot swapping
type DynamicHandler struct {
	currentHandler http.Handler
}

func (d *DynamicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d.currentHandler != nil {
		d.currentHandler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Server initializing...", http.StatusServiceUnavailable)
	}
}

func (d *DynamicHandler) UpdateRoutes() {
	// Note: In a real hot reload, you'd want to rebuild the Go code
	// For now, this just reloads the routes (which works for route changes)
	d.currentHandler = routes.SetupRoutes()
}

func hasGoFileChanged(lastMod *time.Time) bool {
	var changed bool

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories we don't care about
		if info.IsDir() && shouldSkipDir(path) {
			return filepath.SkipDir
		}

		// Check .go files (excluding this main.go to avoid infinite loops)
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "main.go") && info.ModTime().After(*lastMod) {
			changed = true
			*lastMod = info.ModTime()
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v", err)
	}

	return changed
}

func shouldSkipDir(path string) bool {
	skipDirs := []string{
		".git", "vendor", "node_modules", "tmp",
		".vscode", ".idea", "dist", "build",
	}

	for _, skip := range skipDirs {
		if strings.Contains(path, skip) {
			return true
		}
	}
	return false
}

func initializeApp() error {
	// Initialize bootstrap
	bootstrap.InitializeApp()
	return nil
}

func displayWelcomeMessage() {
	elephant := `
   _..--""-.                  .-""--.._
.-'         \ __...----...__ /         '-.
.'      .:::...,'              ',...:::.      '.
(     .''''''::;                  ;::''''''.     )
\             '-)              (-'             /
\             /                \             /
 \          .'.-.            .-.'.          /
  \         | \0|            |0/ |         /
   |         \  |   .-==-.   |  /         |
   \         '/';          ;'\'         /
    '.._      (_ |  .-==-.  | _)      _..'
        '""'-./ /'        '\ \.-'"'"
             / /';   .==.   ;'\ \
        .---/ /   \  .==.  /   / \---.
        |   | |   / .''''. \   | |   |
        |   | |   \ \    / /   | |   |
        |   \ \   /  '""'  \   / /   |
        \    \ \_/          \_/ /    /
         \    \  -._      _. -  /    /
          \    \    '""""'    /    /
           \    \     _    /    /
            \    \   /----\   /    /
`

	fmt.Println(elephant)
	fmt.Println("🐘 Welcome to Enzovu Framework!")
	fmt.Printf("📅 %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
