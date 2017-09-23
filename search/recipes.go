package search

import (
	"context"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

// RecipeType that denotes one recipe
var RecipeType = "recipe"

// RecipeMapping for index
var RecipeMapping = fmt.Sprintf(`{
	"%s": {
		"properties": {
			"title": {
               "type": "keyword",
			   "fields": {
	              "lowercase": {
	                 "type": "string",
	                 "analyzer": "custom_autocomplete"
	              }
	           }
            }
		}
	}
}`, RecipeType)

// ByIngredient search for recipes matching the terms
func ByIngredient(values []string, from int, size int) (*elastic.SearchResult, error) {
	ctx := context.Background()
	client, err := elastic.NewClient()
	// client, err := elastic.NewClient(
	// 	elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	// 	elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	// 	elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]]", 0)))

	if err != nil {
		return nil, err
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
		Type(RecipeType).
		Query(query).
		Pretty(true).
		Do(ctx)

	return response, err
}
