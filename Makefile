.DEFAULT_GOAL := build

update:
	go get -u
	go mod tidy

build:
	go fmt ./...
	go vet ./...
	go build

lint:
	# wget https://github.com/golangci/golangci-lint/releases/download/v1.24.0/golangci-lint-1.24.0-linux-amd64.tar.gz
	./golangci-lint run ./...
	go mod tidy

test: build
	go test -race ./...
