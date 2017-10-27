package search

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Index for storing cocktail recipes and categories
var Index = "cocktails"

// Meta data including pagination
type Meta struct {
	Pagination struct {
		Total       int `json:"total"`
		Count       int `json:"count"`
		PerPage     int `json:"per_page"`
		CurrentPage int `json:"current_page"`
		TotalPages  int `json:"total_pages"`
		Links       struct {
			Next string `json:"next"`
		} `json:"links"`
	} `json:"pagination"`
}

// Mapping for index
var Mapping = `
{
	"settings":{
		"analysis":{
			"filter": {
		        "autocomplete_filter": {
		            "type":     "edge_ngram",
		            "min_gram": 1,
		            "max_gram": 20
		        }
		    },
			"analyzer":{
				"custom_lower":{
					"type":"custom",
					"tokenizer":"lowercase"
				},
				"custom_autocomplete": {
		            "type":      "custom",
		            "tokenizer": "standard",
		            "filter": [
		                "lowercase",
		                "autocomplete_filter"
		            ]
		        }
			}
		}
	}
}
`

// CreateIndex with given name
func CreateIndex(index string) error {
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
		_, err = client.DeleteIndex(index).Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	_, err = client.CreateIndex(index).BodyString(Mapping).Do(ctx)
	if err != nil {
		panic(err)
	}

	return err
}

// CreateMapping of type
func CreateMapping(index string, tp string, mapping string) error {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return err
	}
	_, err = client.PutMapping().Index(index).Type(tp).BodyString(mapping).Do(ctx)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}

	return nil
}

// Save items into a new index with a type
func Save(items []interface{}, index string, tp string) error {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	bulkRequest := client.Bulk()

	for _, item := range items {
		id := formatID(item)

		indexReq := elastic.
			NewBulkIndexRequest().
			Index(index).
			Type(tp).
			Id(id).
			Doc(item)
		bulkRequest = bulkRequest.Add(indexReq)
	}

	res, err := bulkRequest.Do(ctx)

	if err != nil {
		return err
	}

	fmt.Printf("\nIndexed %d/%d items: ", len(res.Succeeded()), len(items))
	for _, item := range res.Succeeded() {
		fmt.Printf(item.Id + " ")
	}
	fmt.Println("")

	return nil
}

// Get returns an item from an index by id
func Get(id string, typ string, index string) (*elastic.GetResult, error) {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Get().Index(index).Type(typ).Id(id).Do(ctx)
	if err != nil || response.Found == false {
		return nil, err
	}

	return response, err
}

// Find items matching all criteria
func Find(terms map[string]string, typ string, index string, size int, from int, sortField string, asc bool) (*elastic.SearchResult, error) {
	query := elastic.NewBoolQuery()

	for field, term := range terms {
		if len(term) > 0 {
			q := elastic.NewMatchQuery(field, term).Operator("AND")
			query = query.Must(q)
		}
	}

	return FindByQuery(query, typ, index, size, from, sortField, asc)
}

// FindAny items matching any of the criteria
func FindAny(terms map[string]string, typ string, index string, size int, from int, sortField string, asc bool) (*elastic.SearchResult, error) {
	query := elastic.NewBoolQuery()

	for field, term := range terms {
		if len(term) > 0 {
			q := elastic.NewMatchQuery(field, term).Operator("OR")
			query = query.Should(q)
		}
	}

	return FindByQuery(query, typ, index, size, from, sortField, asc)
}

// FindByQuery search for a result by a term
func FindByQuery(query elastic.Query, typ string, index string, size int, from int, sortField string, asc bool) (*elastic.SearchResult, error) {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return nil, err
	}

	response, err := client.
		Search(index).
		From(from).
		Size(size).
		Type(typ).
		Pretty(true).
		Sort(sortField, asc).
		Query(query).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return response, err
}

func formatID(item interface{}) string {
	itemID := item.(map[string]interface{})["id"]
	typ := reflect.TypeOf(itemID).String()
	var id string

	if typ == "string" {
		id = itemID.(string)
	} else if typ == "float64" {
		id = strconv.Itoa(int(itemID.(float64)))
	} else {
		panic(fmt.Sprintf("Invalid index ID format: '%v'", typ))
	}

	return id
}
