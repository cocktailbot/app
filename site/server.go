package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/search"
	"github.com/cocktailbot/app/site/controllers"
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

	recipes := new(controllers.Recipes)
	home := new(controllers.Home)

	http.HandleFunc(controllers.RecipesDetailPath, recipes.Detail)
	http.HandleFunc(controllers.RecipesSearchPath, recipes.Search)
	http.HandleFunc(controllers.HomePath, home.Index)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
