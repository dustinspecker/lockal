.PHONY: build
build:
	go build -o ./bin/lockal ./cmd/lockal/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test-integration
test-integration:
	./test/integration/examples.sh

.PHONY: test-unit
test-unit:
	go test ./... -cover -coverprofile=cover.out -covermode=count
