build-app: echo "build-app"

test-app: exit 0


GRPC_PROTO_PATH = ./api/worker-handler.proto
GRPC_GEN_PATH = .

grpc-generate:
	protoc \
		--go_out=${GRPC_GEN_PATH}  --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_GEN_PATH) --go-grpc_opt=paths=source_relative \
    ${GRPC_PROTO_PATH}

make: kubectl apply -f build/deploy.yml

delete: kubectl delete deployment worker-deployment
