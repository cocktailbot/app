package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Recipes represents a collection of recipes
type Recipes struct {
	Data []struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		Categories []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"categories"`
		DifficultyRating string `json:"difficultyRating"`
		RecipeTimes      []struct {
			Title string `json:"title"`
			Time  string `json:"time"`
		} `json:"recipeTimes"`
		TotalTime   string `json:"totalTime"`
		Description string `json:"description"`
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
	} `json:"data"`
}

// Categories are a recipe taxonomy
type Categories struct {
	Data []struct {
		ID       string `json:"id"`
		Slug     string `json:"slug"`
		Title    string `json:"title"`
		Children []struct {
			ID    string `json:"id"`
			Slug  string `json:"slug"`
			Title string `json:"title"`
		} `json:"children,omitempty"`
	} `json:"data"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Get path arg
	args := os.Args[1:]
	path := args[0]
	fmt.Println("Reading file: " + path)

	// Read in json string
	data, err := ioutil.ReadFile(path)
	check(err)
	fmt.Println(string(data))

	// Parse json to struct
	var dat map[string]interface{}

	if err := json.Unmarshal(data, &dat); err != nil {
		panic(err)
	}

	fmt.Println(dat)
}
