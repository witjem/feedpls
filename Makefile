tests:
	go test ./...

lint:
	golangci-lint run
