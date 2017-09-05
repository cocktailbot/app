package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
)

const prefix = "./resources/"

func pageHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	serveTemplate(w, r, "index.html", nil)
}

func pageSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	data := map[string]interface{}{
		"Test":        "oh hello",
		"Results":     results,
		"Ingredients": strings.Join(query["ingredients"], ","),
	}

	serveTemplate(w, r, "search.html", data)
}

func pageRecipesDetail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/recipes/"):]
	recipe, e := search.Get(id)
	err.Check(e)
	data := map[string]interface{}{
		"Recipe": recipe,
	}

	serveTemplate(w, r, "recipes/index.html", data)
}

func apiRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/recipes", apiRecipes)

	fs := http.FileServer(http.Dir(prefix + "static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/recipes/", pageSearch)
	http.HandleFunc("/search", pageSearch)
	http.HandleFunc("/", pageHome)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request, path string, data interface{}) {

	lp := filepath.Join(prefix, "templates", "layout.html")
	fp := filepath.Join(prefix, "templates", "/pages/"+filepath.Clean(path))

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
