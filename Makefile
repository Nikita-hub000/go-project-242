build:
	go build ./cmd/hexlet-path-size

test:
	go test ./...

lint:
	golangci-lint run
	golangci-lint fmt --diff

lint-fix:
	golangci-lint run --fix
	golangci-lint fmt
