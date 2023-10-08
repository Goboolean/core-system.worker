SQL_API_PATH=./api/sql

proto-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
    ./api/kafka/model.*/event.proto

make-app: kubectl apply -f build/deploy.yml

delete: kubectl delete deployment worker-deployment

build-dockerfile: docker build -t worker -f build/Dockerfile .

get-schema:
	curl -L https://raw.githubusercontent.com/Goboolean/shared/main/api/sql/schema.sql -o ${SQL_API_PATH}/schema.sql

sqlc-generate:
	sqlc generate -f ${SQL_API_PATH}/sqlc.yml
	
sqlc-check:
	sqlc compile -f ${SQL_API_PATH}/sqlc.yml

