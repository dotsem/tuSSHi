lint:
    golangci-lint run

fmt:
    gofmt -w -s .

run:
    go run cmd/tusshi/main.go