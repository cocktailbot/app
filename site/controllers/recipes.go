package controllers

import (
	"net/http"
	"strings"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/search"
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
func (c Recipes) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ingredients := c.Param("ingredients", query, "")
	size := c.ParamInt("per_page", query, 10)
	page := c.ParamInt("page", query, 1)
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

	c.Render(w, r, "search.html", data)
}

// Detail page for one recipe
func (c Recipes) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(RecipesDetailPath):]
	id := strings.Split(slug, "-")[0]
	recipe, e := search.Get(id)
	err.Check(e)

	if recipe.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Recipe": recipe,
	}

	c.Render(w, r, "/recipes/index.html", data)
}
