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
	size := 10000 // Target all categories
	page := 0

	results := []models.Category{}
	response, e := search.FindAll(size, page, search.CategoryType, search.Index, "title", true)
	err.Check(e)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var c models.Category
			e = json.Unmarshal(*hit.Source, &c)
			err.Check(e)
			results = append(results, c)
		}
	}
	data := map[string]interface{}{
		"Categories": results,
	}

	c.Render(w, r, "categories/index.html", data)
}

// Detail page for a category
func (c Categories) Detail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(CategoriesDetailPath):]
	slugs := strings.Split(path, "/")
	slug := slugs[len(slugs)-1]
	category := new(models.Category)
	values := map[string]string{
		"slug":                   slug,
		"children.slug":          slug,
		"children.children.slug": slug,
	}
	response, e := search.GetBy(values, search.CategoryType, search.Index)

	if response == nil || response.TotalHits() != 1 {
		http.NotFound(w, r)
		return
	}
	err.Check(e)
	e = json.Unmarshal(*response.Hits.Hits[0].Source, &category)
	err.Check(e)

	recipes := []models.Recipe{}
	values = map[string]string{
		"categories.slug": slug,
	}
	response, e = search.GetBy(values, search.RecipeType, search.Index)
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
