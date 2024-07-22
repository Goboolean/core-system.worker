SQL_API_PATH=./api/sql
.PHONY: test
WIRE_PACKAGES = ./internal/job/fetcher ./internal/job/executer ./internal/job/transmitter/v1 

proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
    ./api/kafka/model.*/event.proto

make-app: 
	kubectl apply -f build/deploy.yml

delete: 
	kubectl delete deployment worker-deployment

build-dockerfile: 
	docker build -t worker -f build/Dockerfile .

wire-job: 
	wire ${WIRE_PACKAGES}

test:
	go clean -testcache
	go test ./... -tags develop

 