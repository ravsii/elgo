test:
	go test ./... -timeout 1m --race

coverage:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

lint:
	golangci-lint run ./... --out-format colored-line-number --config ./.golangci.yml --fix

generate:
	protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    grpc/pb/pool.proto