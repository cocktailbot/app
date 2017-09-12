package main

import (
	"log"
	"net/http"

	"github.com/cocktailbot/app/site/controllers"
)

func main() {
	fs := http.FileServer(http.Dir(controllers.Prefix + "static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	recipes := new(controllers.Recipes)
	categories := new(controllers.Categories)
	home := new(controllers.Home)

	http.HandleFunc(controllers.RecipesDetailPath, recipes.Detail)
	http.HandleFunc(controllers.RecipesSearchPath, recipes.Search)
	http.HandleFunc(controllers.CategoriesDetailPath, categories.Detail)
	http.HandleFunc(controllers.HomePath, home.Index)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
