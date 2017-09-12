package search

import (
	"context"
	"encoding/json"
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

// BaseIndexable common methods and properties
type BaseIndexable struct {
	ID string `json:"id"`
}

// GetID returns unique id
func (b BaseIndexable) GetID() (id string) {
	return b.ID
}

// BaseCollection common methods and properties
type BaseCollection struct {
	Data []Indexable `json:"data"`
}

// GetData returns collection
func (bc BaseCollection) GetData() []Indexable {
	return bc.Data
}

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

	for i := 0; i < len(items.GetData()); i++ {
		item := items.GetData()[i]
		indexReq := elastic.
			NewBulkIndexRequest().
			Index(index).
			Type(tp).
			Id(string(item.GetID())).
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
func Get(id string, index string, item Indexable) (err error) {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	response, err := client.Get().Index(index).Id(id).Do(ctx)
	if err != nil || response.Found == false {
		return err
	}

	err = json.Unmarshal(*response.Source, item)
	return err
}
