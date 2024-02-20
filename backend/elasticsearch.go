package backend

import (
	"context"
	"fmt"
	"log"

	"github.com/dp3why/mrgo/constants"

	"github.com/olivere/elastic/v7"
)

var _ ElasticsearchService = &ElasticsearchBackend{}

var (
    ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
    client *elastic.Client
}

func InitElasticsearchBackend() {
	// connection
    client, err := elastic.NewClient(
		elastic.SetSniff(false), 
        elastic.SetURL(constants.ES_AWS_URL),
        elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_AWS_PASSWORD))
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }

	// part1 post index
    exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
    if err != nil {
        log.Fatalf("Error checking if the index exists: %s", err)
    }

    if !exists {
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            log.Fatalf("Error creating the post index: %s", err)
        }
    }


	// part2 user index
    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            log.Fatalf("Error creating the user index: %s", err)
        }
    }
    fmt.Println("Indexes are created.")

    ESBackend = &ElasticsearchBackend{
        client: client,
    }
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    if err != nil {
        return err
    }
    return nil
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return err
    }
    return nil
}
