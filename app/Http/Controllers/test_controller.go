package controllers

import "net/http"

// Test handles HTTP requests for Test
func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Test controller!"))
}
