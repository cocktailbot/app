package controllers

import "net/http"

// HomePath path to homepage
const HomePath = "/"

// Home controller
type Home struct {
	Application
}

// Index for the home page
func (controller Home) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != HomePath {
		http.NotFound(w, r)
		return
	}

	controller.Render(w, r, "home/index.html", nil)
}
