# godan-dist
WIP

## godanserver backends
Right now godanserver support the following backends:
* MySQL
* ElasticSearch

To switch between backends you need to modify the godanserver config file (godanserver.toml) and set the right Type value 
### MySQL config
```
[DB]
Type = "mysql"
Host = "db"
Port = "3306"
Username = "godan"
Password = "change_this_pwd!"
Name = "godan"
```
### ElasticSearch config
```
[DB]
Type = "elasticsearch"
Host = "db"
Port = "9200"
Username = "godan"
Password = "change_this_pwd!"
Name = "godan"
```
Note: ElasticSearch currently ignores Username & Password parameters
