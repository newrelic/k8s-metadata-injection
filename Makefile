OSFLAG := $(shell uname -s | tr A-Z a-z)
OSFLAG := $(OSFLAG)_amd64
BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME = k8s-metadata-injection

GOVENDOR_VERSION = 1.0.8
GOLANGCILINT_VERSION = 1.12

# required for enabling Go modules
export GO111MODULE=on

.PHONY: all
all: build

.PHONY: build
build: clean lint test compile # TODO: add license-check

.PHONY: clean
clean:
	@echo "[clean] Removing binaries"
	@rm -rf $(BIN_DIR)/$(BINARY_NAME)

$(TOOLS_DIR):
	@mkdir -p $@

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	@echo "[tools] Downloading 'golangci-lint'"
	@wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINDIR=$(@D) sh -s v$(GOLANGCILINT_VERSION) &> /dev/null

# $(TOOLS_DIR)/papers-go: $(TOOLS_DIR)
# 	@echo "[tools] Downloading 'papers-go'"
# 	@go get go.datanerd.us/p/ohai/papers-go/...
# 	@cp $(GOPATH)/bin/papers-go $(TOOLS_DIR)/papers-go

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint
	@echo "[validate] Validating source code running golangci-lint"
	@$(TOOLS_DIR)/golangci-lint run

.PHONY: lint-all
lint-all: $(TOOLS_DIR)/golangci-lint
	@echo "[validate] Validating source code running golangci-lint"
	@$(TOOLS_DIR)/golangci-lint run --enable=interfacer --enable=gosimple

# .PHONY: license-check
# license-check: $(TOOLS_DIR)/papers-go
# 	@echo "[validate] Validating licenses of package dependencies required by the project"
# 	@$(TOOLS_DIR)/papers-go validate

.PHONY: compile
compile:
	@echo "[compile] Building $(BINARY_NAME)"
	@go build -o $(BIN_DIR)/$(BINARY_NAME)

.PHONY: compile-dev
compile-dev:
	@echo "[compile-dev] Building $(BINARY_NAME) for development environment"
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)

.PHONY: deploy-dev
deploy-dev: compile-dev
	@echo "[deploy-dev] Deploying dev container image containing $(BINARY_NAME) in Kubernetes"
	@skaffold run

.PHONY: test
test:
	@echo "[test] Running unit tests"
	@go test .

.PHONY: benchmark-test
benchmark-test:
	@echo "[test] Running benchmark tests"
	@go test -run=^Benchmark* -bench .