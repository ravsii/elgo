test:
	go test ./... -timeout 1m --race

coverage:
	go test ./... -cover -timeout 1m --race

lint:
	golangci-lint run ./... --out-format colored-line-number --config ./.golangci.yml --fix