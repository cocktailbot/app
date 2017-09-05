package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Application container
type Application struct {
}

// Prefix path for template location
const Prefix = "./resources/"

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

	lp := filepath.Join(Prefix, "templates", "layout.html")
	fp := filepath.Join(Prefix, "templates", "/pages/"+filepath.Clean(path))

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

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

// Param return query parameter or default
func (a Application) Param(param string, r *http.Request, def string) (value string) {
	query := r.URL.Query()
	if query[param] == nil {
		query[param] = []string{def}
	}

	return strings.Join(query[param], ",")
}

// ParamInt something
func (a Application) ParamInt(param string, r *http.Request, def int) (value int) {
	str := a.Param(param, r, string(def))
	i, e := strconv.ParseInt(str, 10, 32)
	if e != nil {
		return def
	}

	return int(i)
}
