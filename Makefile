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
	cd docker/server/ && docker build -t zlowram/godanserver .

build-worker:
	cd docker/worker/ && docker build -t zlowram/godanworker .

build-ui:
	cd docker/ui/api/ && docker build -t zlowram/godanapiui .

build-web-ui:
	cd docker/ui/webui/ && docker build -t zlowram/godanwebui .

create:
	cd docker && ./run_containers.sh 

destroy:
	docker kill db rabbitmq godan_server godan_worker godan_api_ui ui_db && docker rm db rabbitmq godan_server godan_worker godan_api_ui ui_db
