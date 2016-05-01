#!/bin/bash

docker run --name db -e MYSQL_ROOT_PASSWORD=change_this_pwd! -e MYSQL_USER=godan -e MYSQL_PASSWORD=change_this_pwd! -e MYSQL_DATABASE=godan -d mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
docker run --name rabbitmq -p 15672:15672 -p 5672:5672 -d rabbitmq:management
docker run --name godan_worker --link rabbitmq -d zlowram/godanworker
docker run --name godan_server -p 8080:8080 --link rabbitmq --link db -d zlowram/godanserver

docker run --name ui_db -p 27017:27017 -d mongo:3.0.4
docker run --name godan_api_ui -p 8000:8000 --link ui_db --link godan_server -d zlowram/godanapiui
docker run --name godan_web_ui -p 8081:80 --link godan_api_ui -d zlowram/godanwebui
