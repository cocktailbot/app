package controllers

import (
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

// Detail page for one recipe
func (c Categories) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(CategoriesDetailPath):]
	category := new(search.Category)
	e := search.OneBy(slug, "slug", category)
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
