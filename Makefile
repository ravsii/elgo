test:
	go test ./... -timeout 10s --count 1 --race

coverage:
	go test ./... -cover -timeout 1m --count 1 --race

lint:
	golangci-lint run ./... --out-format colored-line-number --config ./.golangci.yml --fix