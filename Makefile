test:
	go test ./... -timeout 1m --race

coverage:
	go test ./... -cover -timeout 1m --race

lint:
	golangci-lint run ./... --out-format colored-line-number --config ./.golangci.yml --fix

generate:
	protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    grpc/schema/pool.proto