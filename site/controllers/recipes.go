package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/models"
	"github.com/cocktailbot/app/search"
)

const (
	// RecipesSearchPath points to search results page
	RecipesSearchPath = "/search"
	// RecipesIndexPath list page
	RecipesIndexPath = "/recipes"
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

	results := []models.Recipe{}
	response, e := search.ByIngredient(strings.Split(ingredients, ","), int(page), int(size))
	err.Check(e)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var r models.Recipe
			e = json.Unmarshal(*hit.Source, &r)
			err.Check(e)
			results = append(results, r)
		}
	}
	data := map[string]interface{}{
		"Results":     results,
		"Ingredients": ingredients,
	}

	c.Render(w, r, "recipes/search.html", data)
}

// Detail page for one recipe
func (c Recipes) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(RecipesDetailPath):]
	id := strings.Split(slug, "-")[0]
	recipe := new(models.Recipe)
	response, e := search.Get(id, search.RecipeType, search.Index)
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

	c.Render(w, r, "recipes/detail.html", data)
}

// Index page
func (c Recipes) Index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	terms := map[string]string{
		"ingredients.list.ingredient": strings.ToLower(query.Get("ingredients")),
		"title.lowercase":             strings.ToLower(query.Get("title")),
	}

	size := c.ParamInt("per_page", query, 10000)
	page := c.ParamInt("page", query, 1)
	page = (page - 1) * size
	if page < 0 {
		page = 0
	}

	results := []models.Recipe{}
	response, e := search.Find(terms, size, page, search.RecipeType, search.Index, "title", true)
	err.Check(e)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var c models.Recipe
			e = json.Unmarshal(*hit.Source, &c)
			err.Check(e)
			results = append(results, c)
		}
	}

	pagination := createPagination(page, size, int(response.Hits.TotalHits))
	data := map[string]interface{}{
		"Recipes":    results,
		"Pagination": pagination,
	}

	c.Render(w, r, "recipes/index.html", data)
}
