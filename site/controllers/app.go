package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Application container
type Application struct {
}

// Prefix path for template location
const Prefix = "./resources/"

// Render tries to write html template, or throw 404 if not found
func (c Application) Render(w http.ResponseWriter, r *http.Request, path string, data interface{}) {

	lp := filepath.Join(Prefix, "templates", "layout.html")
	fp := filepath.Join(Prefix, "templates", "/pages/"+filepath.Clean(path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
