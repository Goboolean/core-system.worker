proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
    ./api/model.*/event.proto

make-app: kubectl apply -f build/deploy.yml

delete: kubectl delete deployment worker-deployment

build-dockerfile: docker build -t worker -f build/Dockerfile .