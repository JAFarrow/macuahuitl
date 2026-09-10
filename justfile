dev:
    go run .

build:
    cd frontend && npm run build
    go build -o macuahuitl .

lint:
    go vet ./...
    golangci-lint run
