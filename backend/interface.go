package backend

import (
	"github.com/olivere/elastic/v7"
)

// ElasticsearchService defines the interface for interacting with Elasticsearch.
type ElasticsearchService interface {
	ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error)
	SaveToES(i interface{}, index string, id string) error
	DeleteFromES(query elastic.Query, index string) error
}
