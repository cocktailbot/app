package recipes

import (
	"context"
	"encoding/json"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Find something
func Find() {
	ctx := context.Background()
	client, err := elastic.NewClient()
	check(err)

	index := "cocktails"
	// typ := "recipe"

	//
	// get1, err := client.Get().
	// 	Index(index).
	// 	Type(typ).
	// 	Id(recipes.Data[4].ID).
	// 	Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	//
	var recipe Recipe
	// if get1.Found {
	// 	data, _ := get1.Source.MarshalJSON()
	// 	json.Unmarshal(data, &recipe)
	// 	// fmt.Println(recipe.Title)
	// 	// fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	// }

	// Search with a term query
	query := elastic.NewTermQuery("id", 861)
	// query := elastic.NewFuzzyQuery("title", recipes.Data[0].Title)
	// query := elastic.NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(index).
		Query(query).
		// Sort("user", true). // sort by "user" field, ascending
		// From(0).Size(10).
		// Pretty(true).
		Do(ctx)
	check(err)

	fmt.Printf("Results %d\n", len(searchResult.Hits.Hits))

	for i := 0; i < int(searchResult.Hits.TotalHits); i++ {
		// fmt.Println(searchResult.Hits.Hits[i].Index)
		// fmt.Println(searchResult.Hits.Hits[i].Type)
		// fmt.Println(searchResult.Hits.Hits[i].Id)
		data, _ := searchResult.Hits.Hits[i].Source.MarshalJSON()
		json.Unmarshal(data, &recipe)
		fmt.Println(recipe.Title)
	}

	// data, _ := searchResult.Hits.Hits[0].Source.MarshalJSON()
	// json.Unmarshal(data, &recipe)
	// fmt.Println(recipe.ID)
	// fmt.Println(searchResult.Hits.TotalHits)
	// // searchResult is of type SearchResult and returns hits, suggestions,
	// // and all kinds of other information from Elasticsearch.
	// fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
}
