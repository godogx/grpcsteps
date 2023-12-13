MODULE_NAME=grpcsteps

VENDOR_DIR = vendor

GOLANGCI_LINT_VERSION ?= v1.55.2

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
GHERKIN_LINT ?= gherkin-lint

GITHUB_OUTPUT ?= /dev/null

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: lint
lint: lint-go lint-gherkin

.PHONY: lint-go
lint-go: $(GOLANGCI_LINT) $(VENDOR_DIR)
	@$(GOLANGCI_LINT) run -c .golangci.yaml

.PHONY: lint-gherkin
lint-gherkin:
	@$(GHERKIN_LINT) -c .gherkin-lintrc features/*

.PHONY: test
test: test-unit

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

.PHONY: gen
gen:
	@rm -rf internal/grpctest
	@protoc --go_out=. --go-grpc_out=. resources/protobuf/service.proto

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "MODULE_NAME=$(MODULE_NAME)" >> "$@"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >> "$@"

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)
