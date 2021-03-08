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
	"$$(go env GOPATH)/bin/golangci-lint" run ./...
	go mod tidy

.PHONY: lint-update
lint-update:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint --version

.PHONY: lint-docker
lint-docker:
	docker pull golangci/golangci-lint:latest
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

.PHONY: test
test: build
	go test -race ./...
