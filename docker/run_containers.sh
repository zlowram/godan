#!/bin/bash
if [ "$1" == "es" ] || [ "$1" == "mysql" ]
then
	cp server/godanserver.toml."${1}" server/godanserver.toml
	cp docker-compose-"${1}".yml docker-compose.yml
	docker-compose build && docker-compose up -d
else
	echo "Usage: $0 es|mysql"
fi
