package models

import "html/template"

// Recipe represents cocktail recipe
type Recipe struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Categories []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
		URL   string `json:"url"`
	} `json:"categories"`
	DifficultyRating string `json:"difficultyRating"`
	RecipeTimes      []struct {
		Title string `json:"title"`
		Time  string `json:"time"`
	} `json:"recipeTimes"`
	TotalTime   string        `json:"totalTime"`
	Serves      string        `json:"serves"`
	Description template.HTML `json:"description"`
	Ingredients []struct {
		Title string `json:"title"`
		List  []struct {
			Amount     string `json:"amount"`
			Ingredient string `json:"ingredient"`
			Notes      string `json:"notes"`
		} `json:"list"`
	} `json:"ingredients"`
	Methods []struct {
		Title string `json:"title"`
		List  []struct {
			Step string `json:"step"`
		} `json:"list"`
	} `json:"methods"`
	Search string `json:"search"`
}
