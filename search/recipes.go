package search

import (
	"context"
	"encoding/json"

	elastic "gopkg.in/olivere/elastic.v5"
)

// RecipeType that denotes one recipe
var RecipeType = "recipe"

// Recipe represents cocktail recipe
type Recipe struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Categories []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
	} `json:"categories"`
	DifficultyRating string `json:"difficultyRating"`
	RecipeTimes      []struct {
		Title string `json:"title"`
		Time  string `json:"time"`
	} `json:"recipeTimes"`
	TotalTime   string `json:"totalTime"`
	Serves      string `json:"serves"`
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
	Search string `json:"search"`
}

// Recipes represents a collection of recipes
type Recipes struct {
	Data []Recipe `json:"data"`
	Meta Meta     `json:"meta"`
}

// GetData returns collection
func (rs Recipes) GetData() []Recipe {
	return rs.Data
}

// GetID returns unique id
func (r Recipe) GetID() string {
	return r.ID
}

// ByIngredient search for recipes matching the terms
func ByIngredient(values []string, from int, size int) (matches []Recipe, err error) {
	ctx := context.Background()
	client, err := elastic.NewClient()
	// client, err := elastic.NewClient(
	// 	elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	// 	elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	// 	elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]]", 0)))

	if err != nil {
		return matches, err
	}

	query := elastic.NewBoolQuery()

	for i := 0; i < len(values); i++ {
		q := elastic.NewMultiMatchQuery(values[i], "ingredients.*")
		query = query.Should(q)
	}

	response, err := client.
		Search(Index).
		From(from).
		Size(size).
		Query(query).
		Pretty(true).
		Do(ctx)

	if err != nil || response.TotalHits() == 0 {
		return matches, err
	}

	if response.Hits.TotalHits > 0 {
		for _, hit := range response.Hits.Hits {
			var r Recipe
			err := json.Unmarshal(*hit.Source, &r)

			if err != nil {
				return nil, err
			}

			matches = append(matches, r)
		}
	}

	return
}
