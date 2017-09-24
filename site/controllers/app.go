package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Application container
type Application struct {
}

// Pagination stores meta about the current page, total pages, amount per page
type Pagination struct {
	Total    int
	PerPage  int
	Page     int
	Next     int
	Previous int
}

// TemplatePath path for template location
var TemplatePath = filepath.Join(".", "resources", "templates")

// StaticPath location of js/css/etc files
var StaticPath = filepath.Join(".", "resources", "static")

var layout = filepath.Join(TemplatePath, "layouts", "layout.html")
var templates = map[string]*template.Template{
	"categories/index.html": template.Must(
		template.ParseFiles(
			layout,
			filepath.Join(TemplatePath, "pages", "categories", "index.html"),
			filepath.Join(TemplatePath, "pages", "categories", "_category.html"),
		)),
}

// JSON return data as json response
func (a Application) JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Render tries to write html template, or throw 404 if not found
func (a Application) Render(w http.ResponseWriter, r *http.Request, path string, data interface{}) {
	if r.Header.Get("Accept") == "application/json" {
		a.JSON(w, r, data)
		return
	}

	// Some complex templates are declared ahead of time
	if templates[path] != nil {
		if err := templates[path].ExecuteTemplate(w, "layout", data); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fp := filepath.Join(TemplatePath, "/pages/"+filepath.Clean(path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(layout, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Param return query parameter or default
func (a Application) Param(param string, query url.Values, def string) (value string) {
	if query[param] == nil {
		query[param] = []string{def}
	}

	return strings.Join(query[param], ",")
}

// ParamInt something
func (a Application) ParamInt(param string, query url.Values, def int) (value int) {
	str := a.Param(param, query, string(def))
	i, e := strconv.ParseInt(str, 10, 32)
	if e != nil {
		return def
	}

	return int(i)
}

func createPagination(page int, perPage int, total int) (pagination Pagination) {
	pagination.Page = page
	pagination.PerPage = perPage
	pagination.Total = int(math.Ceil(float64(total) / float64(perPage)))

	if pagination.Page > 1 {
		pagination.Previous = page - 1
	}

	if pagination.Page < pagination.Total {
		pagination.Next = page + 1
	}

	return pagination
}
