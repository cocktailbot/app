package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
	"github.com/shrwdflrst/cocktailbot/site/controllers"
)

func apiRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/recipes", apiRecipes)

	fs := http.FileServer(http.Dir(controllers.Prefix + "static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/recipes/", func(w http.ResponseWriter, r *http.Request) {
		controllers.Recipes.Detail(w, r)
	})
	// http.HandleFunc("/search", controllers.Recipes.Search)
	// http.HandleFunc(controllers.HomePath, controllers.Home)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
