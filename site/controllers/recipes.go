package controllers

import (
	"net/http"
	"strings"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
)

const (
	// RecipesSearchPath points to search results page
	RecipesSearchPath = "/search"
	// RecipesDetailPath points to details for a recipe
	RecipesDetailPath = "/recipes/"
)

// Recipes controller
type Recipes struct {
	Application
}

// Search page
func (controller Recipes) Search(w http.ResponseWriter, r *http.Request) {
	ingredients := controller.Param("ingredients", r, "")
	size := controller.ParamInt("page", r, 10)
	page := controller.ParamInt("page", r, 1)
	page = (page - 1) * size

	if page < 0 {
		page = 0
	}

	results, e := search.ByIngredient(strings.Split(ingredients, ","), int(page), int(size))
	err.Check(e)
	data := map[string]interface{}{
		"Results":     results,
		"Ingredients": ingredients,
	}

	controller.Render(w, r, "search.html", data)
}

// Detail page for one recipe
func (controller Recipes) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(RecipesDetailPath):]
	recipe, e := search.Get(id)
	err.Check(e)

	if recipe.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Recipe": recipe,
	}

	controller.Render(w, r, "recipes/index.html", data)
}
