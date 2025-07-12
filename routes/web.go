package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// SetupRoutes configures and returns the main router
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Add logging middleware to all routes
	handler := loggingMiddleware(mux)

	// Static file serving
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// Home route - serves the index.html
	mux.HandleFunc("/", homeHandler)

	// About route
	mux.HandleFunc("/about", aboutHandler)

	// API routes
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/test", testHandler)

	// Test model route (from your existing code)
	mux.HandleFunc("/test-model", testModelHandler)

	return handler
}

// homeHandler serves the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle exact "/" path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Try to serve index.html from public directory
	indexPath := filepath.Join("public", "index.html")
	if fileExists(indexPath) {
		http.ServeFile(w, r, indexPath)
	} else {
		// Fallback to simple welcome message
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Enzovu Framework</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            text-align: center; 
            padding: 50px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            margin: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-direction: column;
        }
        h1 { font-size: 3em; margin-bottom: 20px; }
        p { font-size: 1.2em; margin-bottom: 30px; }
        .links { margin-top: 30px; }
        .links a { 
            color: #ffeb3b; 
            text-decoration: none; 
            margin: 0 15px;
            padding: 10px 20px;
            border: 2px solid #ffeb3b;
            border-radius: 5px;
            transition: all 0.3s;
        }
        .links a:hover { 
            background: #ffeb3b; 
            color: #333; 
        }
    </style>
</head>
<body>
    <h1>🐘 Welcome to Enzovu!</h1>
    <p>Your elegant Go framework is running successfully!</p>
    <div class="links">
        <a href="/about">About</a>
        <a href="/api/health">Health Check</a>
        <a href="/test-model">Test Model</a>
    </div>
</body>
</html>`)
	}
}

// aboutHandler handles the about page
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":   "Welcome to Enzovu Framework",
		"framework": "Enzovu",
		"version":   "1.0.0",
		"language":  "Go",
		"inspired":  "Laravel",
		"features": []string{
			"MVC Architecture",
			"CLI Code Generation",
			"Middleware Support",
			"Elegant Routing",
			"Database Integration",
		},
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler provides a health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "ok",
		"framework": "enzovu",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	}
	json.NewEncoder(w).Encode(response)
}

// testHandler provides a simple test endpoint
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello from Enzovu!",
			"method":  "GET",
		})
	case http.MethodPost:
		json.NewEncoder(w).Encode(map[string]string{
			"message": "POST request received!",
			"method":  "POST",
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
	}
}

// testModelHandler - placeholder for your existing test model functionality
func testModelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Sample user data (replace with your actual model)
	user := map[string]interface{}{
		"id":    1,
		"name":  "Example User",
		"email": "user@example.com",
	}

	json.NewEncoder(w).Encode(user)
}

// loggingMiddleware adds request logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(lrw, r)

		// Log the request
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %d %v\n",
			time.Now().Format("15:04:05"),
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration,
		)
	})
}

// loggingResponseWriter captures the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := http.Dir(".").Open(filename)
	return err == nil
}
