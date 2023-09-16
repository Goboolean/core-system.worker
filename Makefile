GRPC_PROTO_PATH = ./api/worker-handler.proto
GRPC_GEN_PATH = .

grpc-generate:
	protoc \
		--go_out=${GRPC_GEN_PATH}  --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_GEN_PATH) --go-grpc_opt=paths=source_relative \
    ${GRPC_PROTO_PATH}

make-app: kubectl apply -f build/deploy.yml

delete: kubectl delete deployment worker-deployment

build-dockerfile: docker build -t worker -f build/Dockerfile .