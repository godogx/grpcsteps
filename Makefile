VENDOR_DIR = vendor

GO ?= go
GOLANGCI_LINT ?= golangci-lint
GHERKIN_LINT ?= gherkin-lint

.PHONY: $(VENDOR_DIR) lint test test-unit

$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor

lint:
	@$(GOLANGCI_LINT) run
	@$(GHERKIN_LINT) -c .gherkin-lintrc features/*

test: test-unit

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

gen:
	@rm -rf internal/grpctest
	@protoc --go_out=. --go-grpc_out=. resources/protobuf/service.proto
