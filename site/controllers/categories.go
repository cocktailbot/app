package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/search"
)

const (
	// CategoriesSearchPath points to search results page
	CategoriesSearchPath = "/search"
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

	results := []search.Category{}
	response, e := search.FindAll(size, page, search.CategoryType, search.Index)
	err.Check(e)

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var c search.Category
			e = json.Unmarshal(*hit.Source, &c)
			err.Check(e)
			results = append(results, c)
		}
	}
	data := map[string]interface{}{
		"Categories": results,
	}

	c.Render(w, r, "/categories/index.html", data)
}

// Detail page for a category
func (c Categories) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(CategoriesDetailPath):]
	category := new(search.Category)
	// response, e := search.Get(id, search.Index)
	response, e := search.GetBy("slug", slug, search.CategoryType, search.Index)
	fmt.Println(slug)
	fmt.Println(response.TotalHits())
	if response == nil || response.TotalHits() != 1 {
		http.NotFound(w, r)
		return
	}

	err.Check(e)

	e = json.Unmarshal(*response.Hits.Hits[0].Source, &category)
	err.Check(e)

	if category.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Category": category,
	}

	c.Render(w, r, "/categories/detail.html", data)
}
