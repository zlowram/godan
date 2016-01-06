all: compile build compose

compile: compile-server compile-worker

build: build-server build-worker

compile-server:
	cd server && env GOOS=linux GOARCH=amd64 go build && mv server ../docker/server/godanserver

compile-worker:
	cd worker && env GOOS=linux GOARCH=amd64 go build && mv worker ../docker/worker/godanworker

build-server:
	cd docker/server/ && docker build -t zlowram/godanserver .

build-worker:
	cd docker/worker/ && docker build -t zlowram/godanworker .

compose:
	cd docker && docker-compose up

destroy:
	cd docker && docker-compose kill && docker-compose rm -f
