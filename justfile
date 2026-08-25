# Tusshi just file, run via the `just` command runner

# lints the code
lint:
    golangci-lint run

# formats the code
fmt:
    gofmt -w -s .

# runs the unit tests
test:
	go test ./...

# builds the project
build:
    go build -o tusshi cmd/tusshi/main.go

# builds & runs the project. Not intended for prod.
run:
    just build
    ./tusshi

# removes the binary
clean:
    rm ./tusshi

build-release: 
    echo "Todo"