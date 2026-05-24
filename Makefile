build:
	go build ./cmd/hexlet-path-size

test:
	go test ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
