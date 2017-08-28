package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
)

const prefix = "./resources/"

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	serveTemplate(w, r, "index.html")
}

func recipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/recipes", recipes)

	fs := http.FileServer(http.Dir(prefix + "static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request, path string) {
	lp := filepath.Join(prefix, "templates", "layout.html")
	fp := filepath.Join(prefix, "templates", filepath.Clean(path))

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

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
