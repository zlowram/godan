#!/bin/bash

docker run --name db -e MYSQL_ROOT_PASSWORD=change_this_pwd! -e MYSQL_USER=godan -e MYSQL_PASSWORD=change_this_pwd! -e MYSQL_DATABASE=godan -d mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
docker run --name rabbitmq -p 15672:15672 -p 5672:5672 -d rabbitmq:management
docker run --name godan_worker --link rabbitmq -d zlowram/godanworker
docker run --name godan_server -p 8080:8080 --link rabbitmq --link db -d zlowram/godanserver
