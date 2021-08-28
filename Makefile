generate:
	protoc -I . --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	api/service.proto
	protoc -I . --grpc-gateway_out . \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        api/service.proto
lint:
	golangci-lint run
build:
	go test ./...
	docker container rm --force auth 2>/dev/null && docker build -t auth . \
	&& docker run --name auth -e POSTGRES_PASSWORD=somepass -e \
	POSTGRES_USER=postgres -e POSTGRES_DB=postgres --rm -p 6000:5432 -p 8080:8080 -d auth
run:
	docker exec -it auth /auth_bin
