package controllers

import (
	"net/http"
	"strings"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
)

const (
	// SearchPath points to search results page
	RecipesSearchPath = "/search"
	// DetailPath points to details for a recipe
	RecipesDetailPath = "/recipes/"
)

// Recipes controller
type Recipes struct {
	Application
}

// Search page
func (controller Recipes) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	data := map[string]interface{}{
		"Test":        "oh hello",
		"Results":     results,
		"Ingredients": strings.Join(query["ingredients"], ","),
	}

	controller.View(w, r, "search.html", data)
}

// Detail page for one recipe
func (controller Recipes) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/recipes/"):]
	recipe, e := search.Get(id)
	err.Check(e)

	if recipe.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Recipe": recipe,
	}

	controller.View(w, r, "recipes/index.html", data)
}
