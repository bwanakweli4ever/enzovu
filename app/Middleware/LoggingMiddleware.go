package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"enzovu/config"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
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
		logRequest(r, lrw.statusCode, duration)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func logRequest(r *http.Request, statusCode int, duration time.Duration) {
	cfg := config.GetConfig()

	// Choose log format based on environment
	if cfg.App.Environment == "development" {
		// Colorful development logging
		statusColor := getStatusColor(statusCode)
		methodColor := getMethodColor(r.Method)

		fmt.Printf("%s[%s]%s %s%-6s%s %s%-50s%s %s%d%s %s%v%s\n",
			"\033[90m", time.Now().Format("15:04:05"), "\033[0m", // timestamp
			methodColor, r.Method, "\033[0m", // method
			"\033[94m", r.URL.Path, "\033[0m", // path
			statusColor, statusCode, "\033[0m", // status
			"\033[93m", duration, "\033[0m", // duration
		)
	} else {
		// Structured production logging
		log.Printf("method=%s path=%s status=%d duration=%v ip=%s user_agent=%s",
			r.Method,
			r.URL.Path,
			statusCode,
			duration,
			getClientIP(r),
			r.UserAgent(),
		)
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[92m" // green
	case status >= 300 && status < 400:
		return "\033[96m" // cyan
	case status >= 400 && status < 500:
		return "\033[93m" // yellow
	default:
		return "\033[91m" // red
	}
}

func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[94m" // blue
	case "POST":
		return "\033[92m" // green
	case "PUT":
		return "\033[93m" // yellow
	case "DELETE":
		return "\033[91m" // red
	case "PATCH":
		return "\033[95m" // magenta
	default:
		return "\033[90m" // gray
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
