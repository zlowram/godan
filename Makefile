all: compile create

compile: compile-server compile-worker compile-ui copy-web-ui

compile-server:
	export CGO_ENABLED=0 && cd server && env GOOS=linux GOARCH=amd64 go build && mv server ../docker/server/godanserver

compile-worker:
	export CGO_ENABLED=0 && cd worker && env GOOS=linux GOARCH=amd64 go build && mv worker ../docker/worker/godanworker

compile-ui:
	export CGO_ENABLED=0 && cd ui/api && env GOOS=linux GOARCH=amd64 go build && mv api ../../docker/ui/api/godanapiui

copy-web-ui:
	cp -r ui/webui docker/ui/webui/godan_webui

create:
	cd docker && ./run_containers.sh mysql

create-es:
	cd docker && ./run_containers.sh es

destroy:
	cd docker && ./destroy_containers.sh
