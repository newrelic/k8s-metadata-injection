BIN_DIR = ./bin
BINARY_NAME ?= $(BIN_DIR)/k8s-metadata-injection
DOCKER_IMAGE_NAME ?= newrelic/k8s-metadata-injection
# This default tag is used during e2e test execution in the ci
DOCKER_IMAGE_TAG ?= local-dev

GOLANGCILINT_VERSION = 1.43.0

# required for enabling Go modules inside $GOPATH
export GO111MODULE=on

# GOOS and GOARCH will likely come from env
GOOS ?=
GOARCH ?=
CGO_ENABLED ?= 0

ifneq ($(strip $(GOOS)), )
BINARY_NAME := $(BINARY_NAME)-$(GOOS)
endif

ifneq ($(strip $(GOARCH)), )
BINARY_NAME := $(BINARY_NAME)-$(GOARCH)
endif

.PHONY: all
all: build

.PHONY: build
build: test compile

compile:
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(INTEGRATION)..."
	go mod download
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY_NAME) ./cmd/server

.PHONY: compile-multiarch
compile-multiarch:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	$(MAKE) compile GOOS=linux GOARCH=arm64
	$(MAKE) compile GOOS=linux GOARCH=arm

.PHONY: build-container
build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $$DOCKERARGS .

.PHONY: docker-build 
docker-build: 
	docker buildx build --load . -t e2e/k8s-metadata-injection:e2e $$DOCKERARGS .

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

# rt-update-changelog runs the release-toolkit run.sh script by piping it into bash to update the CHANGELOG.md.
# It also passes down to the script all the flags added to the make target. To check all the accepted flags,
# see: https://github.com/newrelic/release-toolkit/blob/main/contrib/ohi-release-notes/run.sh
#  e.g. `make rt-update-changelog -- -v`
rt-update-changelog:
	curl "https://raw.githubusercontent.com/newrelic/release-toolkit/v1/contrib/ohi-release-notes/run.sh" | bash -s -- $(filter-out $@,$(MAKECMDGOALS))


.PHONY: compile rt-update-changelog
