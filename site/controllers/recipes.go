package controllers

import (
	"net/http"
	"strings"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/search"
)

const (
	// CategoriesSearchPath points to search results page
	CategoriesSearchPath = "/search"
	// CategoriesDetailPath points to details for a recipe
	CategoriesDetailPath = "/recipes/"
)

// Categories controller
type Categories struct {
	Application
}

// Detail page for one recipe
func (c Categories) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len(CategoriesDetailPath):]
	id := strings.Split(slug, "-")[0]
	category := new(search.Category)
	e := search.Get(id, search.Index, category)
	err.Check(e)

	if category.ID == "" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Category": category,
	}

	c.Render(w, r, "/categories/index.html", data)
}
