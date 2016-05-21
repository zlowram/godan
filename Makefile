all: compile build create

compile: compile-server compile-worker compile-ui copy-web-ui

build: build-server build-worker build-ui build-web-ui

compile-server:
	export CGO_ENABLED=0 && cd server && env GOOS=linux GOARCH=amd64 go build && mv server ../docker/server/godanserver

compile-worker:
	export CGO_ENABLED=0 && cd worker && env GOOS=linux GOARCH=amd64 go build && mv worker ../docker/worker/godanworker

compile-ui:
	export CGO_ENABLED=0 && cd ui/api && env GOOS=linux GOARCH=amd64 go build && mv api ../../docker/ui/api/godanapiui

copy-web-ui:
	cp -r ui/webui docker/ui/webui/godan_webui

build-server:
	cd docker && docker-compose build godan_server

build-worker:
	cd docker && docker-compose build godan_worker

build-ui:
	cd docker && docker-compose build godan_api_ui

build-web-ui:
	cd docker && docker-compose build godan_web_ui

create:
	cd docker && ./run_containers.sh 

destroy:
	cd docker && ./destroy_containers.sh
