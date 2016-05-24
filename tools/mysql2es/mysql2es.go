package main

import (
	"flag"
	"log"

	"github.com/zlowram/godan/model"
	"github.com/zlowram/godan/persistence"
)

func main() {
	var mysqlDsn = flag.String("mysql", "godan:pwd@tcp(127.0.0.1:3306)/godan?charset=utf8mb4,utf8", "MySQL DSN")
	var elasticUrl = flag.String("es", "http://127.0.0.1:9200", "ElasticSearch URL")
	var elasticIndex = flag.String("index", "godan", "ElasticSearch index name")
	flag.Parse()
	mp := persistence.NewMySQLPersistenceManager(*mysqlDsn)
	defer mp.Close()
	log.Print("Retrieving banners from MySQL database")
	banners, err := mp.QueryBanners(model.Filters{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Retrieved %v banners", len(banners))
	ep := persistence.NewElasticPersistenceManager(*elasticUrl, *elasticIndex)
	for i, banner := range banners {
		if i%1000 == 0 {
			log.Printf("Inserting banner #%v", i)
		}
		ep.SaveBanner(banner)
	}
}
