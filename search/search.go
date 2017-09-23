package search

import (
	"context"
	"fmt"

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

// Collection of indexable items
type Collection interface {
	GetData() []Indexable
}

// Indexable item that has a unique id
type Indexable interface {
	GetID() string
}

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
		client.DeleteIndex(index).Do(ctx)
	}
	_, err = client.CreateIndex(index).Do(ctx)

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
		id := item.(map[string]interface{})["id"].(string)
		indexReq := elastic.
			NewBulkIndexRequest().
			Index(index).
			Type(tp).
			Id(id).
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

// GetBy search for a result by a term
func GetBy(values map[string]string, typ string, index string) (*elastic.SearchResult, error) {
	ctx := context.Background()
	client, err := elastic.NewClient()

	if err != nil {
		return nil, err
	}

	query := elastic.NewBoolQuery()

	for field, term := range values {
		// q := elastic.NewMultiMatchQuery(values[i], "ingredients.*")
		q := elastic.NewTermQuery(field, term)
		query = query.Should(q)
	}

	response, err := client.
		Search(index).
		Type(typ).
		Pretty(true).
		Query(query).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return response, err
}

// FindAll search for a result by a term
func FindAll(size int, from int, typ string, index string, sortField string, asc bool) (*elastic.SearchResult, error) {
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
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return response, err
}
