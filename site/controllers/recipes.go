package controllers

import (
	"encoding/json"
	"fmt"
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

	results := []search.Recipe{}
	response, e := search.ByIngredient(strings.Split(ingredients, ","), int(page), int(size))
	err.Check(e)

	fmt.Printf("%#v", response)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var r search.Recipe
			e = json.Unmarshal(*hit.Source, &r)
			err.Check(e)
			results = append(results, r)
		}
	}
	data := map[string]interface{}{
		"Results":     results,
		"Ingredients": ingredients,
	}

	c.Render(w, r, "/recipes/search.html", data)
}

// Detail page for one recipe
func (c Recipes) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(RecipesDetailPath):]
	id := strings.Split(slug, "-")[0]
	recipe := new(search.Recipe)
	response, e := search.Get(id, search.Index)
	err.Check(e)
	e = json.Unmarshal(*response.Source, &recipe)
	err.Check(e)

	if recipe.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Recipe": recipe,
	}

	c.Render(w, r, "/recipes/detail.html", data)
}
