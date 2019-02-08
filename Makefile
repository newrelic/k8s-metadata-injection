BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME = k8s-metadata-injection
WEBHOOK_DOCKER_IMAGE_NAME=newrelic/k8s-metadata-injection
WEBHOOK_DOCKER_IMAGE_TAG=latest

GOLANGCILINT_VERSION = 1.12

# required for enabling Go modules inside $GOPATH
export GO111MODULE=on

.PHONY: all
all: build

.PHONY: build
build: lint test build-container

$(TOOLS_DIR):
	@mkdir -p $@

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	@echo "[tools] Downloading 'golangci-lint'"
	@wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINDIR=$(@D) sh -s v$(GOLANGCILINT_VERSION) > /dev/null 2>&1

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint
	@echo "[validate] Validating source code running golangci-lint"
	@$(TOOLS_DIR)/golangci-lint run

.PHONY: build-container
build-container:
	docker build -t $(WEBHOOK_DOCKER_IMAGE_NAME):$(WEBHOOK_DOCKER_IMAGE_TAG) .

.PHONY: test
test:
	@echo "[test] Running unit tests"
	@go test ./...

.PHONY: e2e-test
e2e-test:
	@echo "[test] Running e2e tests"
	./e2e-tests/tests.sh

.PHONY: benchmark-test
benchmark-test:
	@echo "[test] Running benchmark tests"
	@go test -run=^Benchmark* -bench .
