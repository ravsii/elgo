test:
	go test ./... -timeout 10s --count 1 --race

coverage:
	go test ./... -cover -timeout 1m --count 1 --race