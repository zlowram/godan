package persistence

import (
	"encoding/json"
	"log"

	"github.com/zlowram/godan/model"

	"gopkg.in/olivere/elastic.v3"
)

type ElasticPersistenceManager struct {
	es    *elastic.Client
	index string
}

func NewElasticPersistenceManager(url, index string) *ElasticPersistenceManager {
	es, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		log.Fatal(err)
	}
	exists, err := es.IndexExists(index).Do()
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		_, err = es.CreateIndex(index).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
	return &ElasticPersistenceManager{es: es, index: index}
}

func (ep *ElasticPersistenceManager) SaveBanner(b model.Banner) {
	_, err := ep.es.Index().Index(ep.index).Type("banner").BodyJson(b).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func (ep *ElasticPersistenceManager) QueryBanners(f model.Filters) ([]model.Banner, error) {
	boolQuery := elastic.NewBoolQuery()
	mustQueries := make([]elastic.Query, 0, 1)
	if f.Ip != "" {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("Ip", f.Ip))
	}
	for _, v := range f.Ports {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("Port", v))
	}
	for _, v := range f.Services {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("Service", v))
	}
	if f.Regexp != "" {
		mustQueries = append(mustQueries, elastic.NewRegexpQuery("Content", f.Regexp))
	}
	if len(mustQueries) > 0 {
		boolQuery.Must(mustQueries...)
	} else {
		boolQuery.Must(elastic.NewMatchAllQuery())
	}

	searchResults, err := ep.es.Search().Index(ep.index).Query(boolQuery).Do()
	if err != nil {
		log.Fatal(err)
	}

	result := make([]model.Banner, 0, searchResults.Hits.TotalHits)
	for _, hit := range searchResults.Hits.Hits {
		var b model.Banner
		err = json.Unmarshal(*hit.Source, &b)
		if err != nil {
			continue
		}
		result = append(result, b)
	}

	return result, nil
}

func (ep *ElasticPersistenceManager) Close() {
}
