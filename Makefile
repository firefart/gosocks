.DEFAULT_GOAL := build

.PHONY: update
update:
	go get -u
	go mod tidy

.PHONY: build
build:
	go fmt ./...
	go vet ./...
	go build

.PHONY: lint
lint:
	@if [ ! -f "$$(go env GOPATH)/bin/golangci-lint" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.34.1; \
	fi
	golangci-lint run ./...
	go mod tidy

.PHONY: test
test: build
	go test -race ./...
