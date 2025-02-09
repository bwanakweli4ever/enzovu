package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Template cache to avoid reloading templates for every request
var templates = make(map[string]*template.Template)

// Render function to process templates and send the response
func Render(w http.ResponseWriter, tmpl string, data interface{}) {
	// Check if the template is already loaded
	if _, ok := templates[tmpl]; !ok {
		// Load the template file from 'resources/views' directory
		tmplPath := filepath.Join("resources", "views", tmpl+".html")
		tmplParsed, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		templates[tmpl] = tmplParsed
	}

	// Render the template with the data
	templates[tmpl].Execute(w, data)
}
