package search

import (
	"context"
	"encoding/json"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Index for storing cocktail recipes
var Index = "cocktails"

// Typ that denotes one recipe
var Typ = "recipe"

// Save items into a new index with a type
func Save(items Collection, index string, tp string) error {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return err
	}

	exists, err := client.IndexExists(index).Do(ctx)

	if err != nil {
		return err
	}
	if exists {
		client.DeleteIndex(index).Do(ctx)
	}

	indexParams := fmt.Sprintf(`{
			"mappings":{
				%s:{
					"properties": {

					}
				}
			}
		}`, tp)

	// Create an index
	_, err = client.CreateIndex(index).BodyString(indexParams).Do(ctx)

	if err != nil {
		return err
	}

	bulkRequest := client.Bulk()

	for i := 0; i < len(items.Data()); i++ {
		item := items.Data()[i]
		indexReq := elastic.
			NewBulkIndexRequest().
			Index(index).
			Type(tp).
			Id(string(item.ID())).
			Doc(item)
		bulkRequest = bulkRequest.Add(indexReq)
	}

	res, err := bulkRequest.Do(ctx)

	fmt.Printf("Indexed items %d\n", len(res.Indexed()))

	if err != nil {
		return err
	}

	return nil
}

// Get returns a recipe by id
func Get(id string) (recipe Recipe, err error) {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return recipe, err
	}

	response, err := client.Get().Index(Index).Id(id).Do(ctx)

	if err != nil || response.Found == false {
		return recipe, err
	}

	err = json.Unmarshal(*response.Source, &recipe)

	return recipe, err
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
