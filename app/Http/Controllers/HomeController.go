package controllers

import (
	"encoding/json"
	models "enzovu/app/Models"
	"enzovu/helpers"
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

type UserController struct{}

// Index handles the user list view
func (u *UserController) RenderView(w http.ResponseWriter, r *http.Request) {
	// Sample data that might come from your model/database
	users := []string{"Alice", "Bob", "Charlie"}
	// Import the debug package from enzovu/helpers
	// Debug the users data
	helpers.Dd(w, users)

	// Render the 'index' view with the users data
	//views.Render(w, "hello", users)
}

func (u *UserController) Show(w http.ResponseWriter, r *http.Request) {
	// Sample user data
	user := map[string]string{
		"name":  "John Doe",
		"email": "john.doe@example.com",
	}

	json.NewEncoder(w).Encode(user)

	// Import the debug package from enzovu/helpers
	// Debug the user data
	// helpers.Dd(w, user)

	// Render the 'show' view with the user data

}
