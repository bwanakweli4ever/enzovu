package main

import (
	"enzovu/bootstrap"
	"enzovu/routes"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting  Enzovu Go Framework...")
	bootstrap.InitializeApp()  // Initialize the app
	routes.RegisterWebRoutes() // Register routes
	http.ListenAndServe(":8000", nil)
}
