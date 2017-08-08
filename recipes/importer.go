package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
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

func main() {
	// Get path arg
	args := os.Args[1:]
	path := args[0]
	fmt.Println("Reading file: " + path)

	// Read in json string
	data, err := ioutil.ReadFile(path)
	check(err)
	// fmt.Println(string(data))

	// Parse json to struct
	var dat map[string]interface{}

	if err = json.Unmarshal(data, &dat); err != nil {
		panic(err)
	}

	// fmt.Println(dat)

	// Create a context
	ctx := context.Background()

	// Create a client
	index := "twitter"
	client, err := elastic.NewClient()
	check(err)

	exists, err := client.IndexExists(index).Do(ctx)
	check(err)
	if exists {
		client.DeleteIndex(index).Do(ctx)
	}

	indexParams := `{
		"mappings":{
			"tweet":{
				"properties": {
					"user": {
						"type":"keyword"
					}
				}
			}
		}
	}`

	// Create an index
	_, err = client.CreateIndex(index).BodyString(indexParams).Do(ctx)
	check(err)

	// Add a document to the index
	tweet := Tweet{User: "olivere", Message: "Take Five"}
	_, err = client.Index().
		Index(index).
		Type("tweet").
		Id("1").
		BodyJson(tweet).
		Refresh("true").
		Do(ctx)
	check(err)

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index(index).       // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	check(err)

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
}
