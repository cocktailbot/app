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
	// CategoriesIndexPath points to categories list
	CategoriesIndexPath = "/categories"
	// CategoriesDetailPath points to details for a recipe
	CategoriesDetailPath = "/categories/"
)

// Categories controller
type Categories struct {
	Application
}

// Index page
func (c Categories) Index(w http.ResponseWriter, r *http.Request) {
	size := 500 // Target all categories
	page := 0
	results := []models.Category{}
	response, e := search.Find(map[string]string{}, search.CategoryType, search.Index, size, page, "title", true)
	err.Check(e)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var c models.Category
			e = json.Unmarshal(*hit.Source, &c)
			err.Check(e)
			results = append(results, c)
		}
	}

	pagination := createPagination(page, size, int(response.Hits.TotalHits))
	data := map[string]interface{}{
		"Categories": results,
		"Pagination": pagination,
	}

	c.Render(w, r, "categories/index.html", data)
}

// Detail page for a category
func (c Categories) Detail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(CategoriesDetailPath):]
	slugs := strings.Split(path, "/")
	slug := slugs[len(slugs)-1]
	category := new(models.Category)
	terms := map[string]string{
		"slug":                   slug,
		"children.slug":          slug,
		"children.children.slug": slug,
	}
	response, e := search.Find(terms, search.CategoryType, search.Index, 1, 0, "title", true)

	if response == nil || response.TotalHits() != 1 {
		http.NotFound(w, r)
		return
	}
	err.Check(e)
	e = json.Unmarshal(*response.Hits.Hits[0].Source, &category)
	err.Check(e)

	recipes := []models.Recipe{}
	terms = map[string]string{
		"categories.slug": slug,
	}
	response, e = search.Find(terms, search.RecipeType, search.Index, 100, 0, "title", true)
	err.Check(e)

	if response != nil && response.TotalHits() > 0 {
		for _, hit := range response.Hits.Hits {
			var r models.Recipe
			e = json.Unmarshal(*hit.Source, &r)
			err.Check(e)
			recipes = append(recipes, r)
		}
	}

	data := map[string]interface{}{
		"Category": category,
		"Recipes":  recipes,
	}

	c.Render(w, r, "/categories/detail.html", data)
}
