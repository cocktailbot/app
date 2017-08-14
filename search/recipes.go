package search

import (
	"context"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Recipe represents cocktail recipe
type Recipe struct {
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
}

// Recipes represents a collection of recipes
type Recipes struct {
	Data []Recipe `json:"data"`
}

// Index for storing cocktail recipes
var Index = "cocktails"

// Typ that denotes one recipe
var Typ = "recipe"

// Save recipes into a new index
func Save(recipes Recipes) error {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return err
	}

	exists, err := client.IndexExists(Index).Do(ctx)

	if err != nil {
		return err
	}
	if exists {
		client.DeleteIndex(Index).Do(ctx)
	}

	indexParams := `{
			"mappings":{
				"recipe":{
					"properties": {

					}
				}
			}
		}`

	// Create an index
	_, err = client.CreateIndex(Index).BodyString(indexParams).Do(ctx)

	if err != nil {
		return err
	}

	bulkRequest := client.Bulk()

	for i := 0; i < len(recipes.Data); i++ {
		fmt.Println(recipes.Data[i].Title)
		indexReq := elastic.
			NewBulkIndexRequest().
			Index(Index).
			Type(Typ).
			Id(recipes.Data[i].ID).
			Doc(recipes.Data[i])
		bulkRequest = bulkRequest.Add(indexReq)
	}

	res, err := bulkRequest.Do(ctx)

	fmt.Printf("Indexed items %d\n", len(res.Indexed()))

	if err != nil {
		return err
	}

	return nil
}
