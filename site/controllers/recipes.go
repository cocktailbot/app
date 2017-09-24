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
	// RecipesIndexPath list page
	RecipesIndexPath = "/recipes"
	// RecipesDetailPath points to details for a recipe
	RecipesDetailPath = "/recipes/"
)

// Recipes controller
type Recipes struct {
	Application
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
	ingredients := strings.ToLower(query.Get("ingredients"))
	title := strings.ToLower(query.Get("title"))
	terms := map[string]string{}

	if len(ingredients) > 0 {
		terms["ingredients.list.ingredient"] = ingredients
	}

	if len(title) > 0 {
		terms["title.lowercase"] = title
	}

	size := c.ParamInt("per_page", query, 10000)
	page := c.ParamInt("page", query, 1)
	page = (page - 1) * size
	if page < 0 {
		page = 0
	}

	results := []models.Recipe{}
	response, e := search.Find(terms, search.RecipeType, search.Index, size, page, "title", true)
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
