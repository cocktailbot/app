package main

import (
	"encoding/json"
	"net/http"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/search"
)

func recipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")
	results, e := search.ByIngredient(query["ingredients"], 0, 10)
	err.Check(e)
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/recipes", recipes)
	http.ListenAndServe(":8080", nil)
}
