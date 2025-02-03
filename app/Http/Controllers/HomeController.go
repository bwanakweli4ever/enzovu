package controllers

import (
	"encoding/json"
	models "enzovu/app/Models"
	"fmt"
	"net/http"
	"path/filepath"
)

// Home handles the request to the homepage and serves the index.html file
func Home(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file from the public directory
	http.ServeFile(w, r, filepath.Join("public", "index.html"))
}

func About(w http.ResponseWriter, r *http.Request) {
	// print a message
	fmt.Fprintf(w, "Welcome to Go-Laravel-like Framework")

}

func TestModel(w http.ResponseWriter, r *http.Request) {
	user := models.GetUser()

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the user data to the response
	json.NewEncoder(w).Encode(user)

}
