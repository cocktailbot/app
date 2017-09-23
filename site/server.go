package main

import (
	"log"
	"net/http"

	"github.com/cocktailbot/app/site/controllers"
)

func main() {
	fs := http.FileServer(http.Dir(controllers.StaticPath + "static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	home := new(controllers.Home)
	recipes := new(controllers.Recipes)
	categories := new(controllers.Categories)

	http.HandleFunc(controllers.HomePath, home.Index)

	http.HandleFunc(controllers.RecipesDetailPath, recipes.Detail)
	http.HandleFunc(controllers.RecipesIndexPath, recipes.Index)
	http.HandleFunc(controllers.RecipesSearchPath, recipes.Search)

	http.HandleFunc(controllers.CategoriesDetailPath, categories.Detail)
	http.HandleFunc(controllers.CategoriesIndexPath, categories.Index)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
