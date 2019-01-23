BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME = k8s-metadata-injection
DOCKER_IMAGE_NAME=quay.io/newrelic/k8s-metadata-injector-dev
DOCKER_IMAGE_TAG=latest

GOLANGCILINT_VERSION = 1.12

# required for enabling Go modules inside $GOPATH
export GO111MODULE=on

.PHONY: all
all: build

.PHONY: build
build: clean lint test docker-build

.PHONY: clean
clean:
	@echo "[clean] Removing binaries"
	@rm -rf $(BIN_DIR)/$(BINARY_NAME)

$(TOOLS_DIR):
	@mkdir -p $@

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	@echo "[tools] Downloading 'golangci-lint'"
	@wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINDIR=$(@D) sh -s v$(GOLANGCILINT_VERSION) &> /dev/null

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint
	@echo "[validate] Validating source code running golangci-lint"
	@$(TOOLS_DIR)/golangci-lint run

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: test
test:
	@echo "[test] Running unit tests"
	@go test .

.PHONY: benchmark-test
benchmark-test:
	@echo "[test] Running benchmark tests"
	@go test -run=^Benchmark* -bench .
