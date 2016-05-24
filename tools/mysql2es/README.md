# mysql2es
Little tool to migrate data from an existing godan MySQL structure to ElasticSearch
Usage:
```
Usage of ./mysql2es:
  -es string
    	ElasticSearch URL (default "http://127.0.0.1:9200")
  -index string
    	ElasticSearch index name (default "godan")
  -mysql string
    	MySQL DSN (default "godan:pwd@tcp(127.0.0.1:3306)/godan?charset=utf8mb4,utf8")
```
