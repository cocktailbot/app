package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

// ByIngredient search for recipes matching the terms
func ByIngredient(values []string, from int, size int) (matches []Recipe, err error) {
	ctx := context.Background()
	// client, err := elastic.NewClient()
	client, err := elastic.NewClient(
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]]", 0)))

	if err != nil {
		return matches, err
	}

	// q := elastic.NewTermQuery("title", values[0])
	// query := elastic.NewTermQuery("id", 861)
	query := elastic.NewBoolQuery()
	// query = query.Should(q)

	for i := 0; i < len(values); i++ {
		q := elastic.NewMultiMatchQuery(values[i], "*")
		fmt.Println(q.Source())
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
