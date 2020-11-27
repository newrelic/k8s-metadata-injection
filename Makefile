BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME ?= k8s-metadata-injection
DOCKER_IMAGE_NAME ?= newrelic/k8s-metadata-injection
DOCKER_IMAGE_TAG ?= 1.3.1

GOLANGCILINT_VERSION = 1.33.0

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
	grep -e "image: $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" deploy/newrelic-metadata-injection.yaml > /dev/null || \
	( echo "Docker image tag being built $(DOCKER_IMAGE_TAG) is not synchronized with deployment yaml" && exit 1 )
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

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

deploy/combined.yaml: deploy/newrelic-metadata-injection.yaml deploy/job.yaml
	echo '---' | cat deploy/newrelic-metadata-injection.yaml - deploy/job.yaml > deploy/combined.yaml
