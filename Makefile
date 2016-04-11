all: compile build create

compile: compile-server compile-worker

build: build-server build-worker

compile-server:
	export CGO_ENABLED=0 && cd server && env GOOS=linux GOARCH=amd64 go build && mv server ../docker/server/godanserver

compile-worker:
	export CGO_ENABLED=0 && cd worker && env GOOS=linux GOARCH=amd64 go build && mv worker ../docker/worker/godanworker

build-server:
	cd docker/server/ && docker build -t zlowram/godanserver .

build-worker:
	cd docker/worker/ && docker build -t zlowram/godanworker .

create:
	cd docker && ./run_containers.sh 

destroy:
	docker kill db rabbitmq godan_server godan_worker && docker rm db rabbitmq godan_server godan_worker
