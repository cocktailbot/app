package recipes

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

// Tweet something
type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Do save json into Elasticsearch
func Do() {
	// Get path arg
	args := os.Args[1:]
	path := args[0]
	// fmt.Println("Reading file: " + path)

	// Read in json string
	file, err := ioutil.ReadFile(path)
	check(err)
	// fmt.Println(string(file))

	// Parse json to struct
	var recipes Recipes

	if err = json.Unmarshal(file, &recipes); err != nil {
		panic(err)
	}

	ctx := context.Background()
	client, err := elastic.NewClient()
	check(err)

	index := "cocktails"
	typ := "recipe"
	exists, err := client.IndexExists(index).Do(ctx)
	check(err)
	if exists {
		client.DeleteIndex(index).Do(ctx)
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
	_, err = client.CreateIndex(index).BodyString(indexParams).Do(ctx)
	check(err)

	bulkRequest := client.Bulk()

	for i := 0; i < len(recipes.Data); i++ {
		// fmt.Println(recipes.Data[i].Title)
		indexReq := elastic.
			NewBulkIndexRequest().
			Index(index).
			Type(typ).
			Id(recipes.Data[i].ID).
			Doc(recipes.Data[i])
		bulkRequest = bulkRequest.Add(indexReq)
	}

	_, err = bulkRequest.Do(ctx)
	check(err)
}
