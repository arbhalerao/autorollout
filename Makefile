IMG ?= autorollout:latest
CONTAINER_TOOL ?= docker
SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: build
build: ## Build the Go binary
	go build -o bin/autorollout ./cmd/main.go

.PHONY: run
run: ## Run the app locally
	go run ./cmd/main.go

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Vet code
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests
	go test ./...

.PHONY: lint
lint: golangci-lint ## Run linter
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run linter with auto-fix
	$(GOLANGCI_LINT) run --fix

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	$(CONTAINER_TOOL) build -t $(IMG) .

.PHONY: docker-push
docker-push: ## Push Docker image
	$(CONTAINER_TOOL) push $(IMG)

.PHONY: docker-buildx
docker-buildx: ## Cross-platform Docker build
	$(CONTAINER_TOOL) buildx build --push --platform linux/amd64,linux/arm64 --tag $(IMG) .

##@ Dev Cluster

.PHONY: create-cluster
create-cluster: ## Create Kind cluster
	@./dev/scripts/create-dev-cluster.sh

.PHONY: delete-cluster
delete-cluster: ## Delete Kind cluster
	@./dev/scripts/delete-dev-cluster.sh

##@ Tools

LOCALBIN ?= $(shell pwd)/bin
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.1.0

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Install golangci-lint if needed
$(GOLANGCI_LINT): $(LOCALBIN)
	@[ -f "$(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION)" ] || { \
		echo "Installing golangci-lint..."; \
		GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
		mv $(GOLANGCI_LINT) $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION); \
	}; \
	ln -sf $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) $(GOLANGCI_LINT)
