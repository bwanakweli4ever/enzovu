package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func Dd(w http.ResponseWriter, value interface{}) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // HTTP 200 OK

	// Get file and line number where Dd() was called
	_, file, line, _ := runtime.Caller(1)

	// Capture stack trace
	stackBuf := make([]byte, 4096)
	stackLen := runtime.Stack(stackBuf, false)

	// Create structured JSON response
	response := map[string]interface{}{
		"value": value,
		"type":  fmt.Sprintf("%T", value),
		"caller": map[string]interface{}{
			"file": file,
			"line": line,
		},
		"stack_trace": string(stackBuf[:stackLen]),
	}

	// Encode response as JSON and send it
	json.NewEncoder(w).Encode(response)
}

func Debug(value interface{}) {
	log.Println(value)
}
