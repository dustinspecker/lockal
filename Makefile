VERSION ?= $(shell git describe --always --dirty --long --tags)

export CGO_ENABLED = 0

.PHONY: build
build:
	go build -ldflags="-X main.VERSION=$(VERSION)" -o ./bin/lockal ./cmd/lockal/main.go

.PHONY: build-linux-amd64
build-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -ldflags="-X main.VERSION=$(VERSION)" -o ./bin/lockal-linux-amd64 ./cmd/lockal/main.go

.PHONY: build-linux-arm64
build-linux-arm64:
	GOARCH=arm64 GOOS=linux go build -ldflags="-X main.VERSION=$(VERSION)" -o ./bin/lockal-linux-arm64 ./cmd/lockal/main.go

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-X main.VERSION=$(VERSION)" -o ./bin/lockal-darwin-amd64 ./cmd/lockal/main.go

.PHONY: cross-build
cross-build: build-linux-amd64 build-linux-arm64 build-darwin-amd64

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test-integration
test-integration:
	./test/integration/examples.sh

.PHONY: test-unit
test-unit:
	go test ./... -cover -coverprofile=cover.out -covermode=count
