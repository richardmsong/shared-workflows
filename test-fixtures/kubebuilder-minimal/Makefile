# Project variables
BINARY_NAME ?= manager
IMG ?= controller:latest
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go variables
GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

# Docker variables
DOCKER_TAGS ?= latest
DOCKER_PUSH ?=
PLATFORMS ?= linux/amd64,linux/arm64
REGISTRY ?= ghcr.io
IMAGE_NAME ?= $(REGISTRY)/richardmsong/shared-workflows-test

# Tools
LOCALBIN ?= $(shell pwd)/bin
KUSTOMIZE ?= $(LOCALBIN)/kustomize
KUSTOMIZE_VERSION ?= v5.3.0

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test -v -race -coverprofile=cover.out ./...

.PHONY: lint
lint: ## Run golangci-lint.
	golangci-lint run

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) .

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: ## Build docker image.
	docker build -t $(IMG) .

.PHONY: docker-push
docker-push: ## Push docker image.
	docker push $(IMG)

.PHONY: docker-buildx-ci
docker-buildx-ci: ## Build and optionally push multi-arch docker image for CI.
	@echo "Building Docker image with tags: $(DOCKER_TAGS)"
	@FULL_TAGS=""; \
	for tag in $(DOCKER_TAGS); do \
		FULL_TAGS="$$FULL_TAGS -t $(IMAGE_NAME):$$tag"; \
	done; \
	if [ -n "$(DOCKER_PUSH)" ]; then \
		echo "Pushing image..."; \
		docker buildx build --platform $(PLATFORMS) $$FULL_TAGS --push .; \
	else \
		echo "Building image (no push)..."; \
		docker buildx build --platform $(PLATFORMS) $$FULL_TAGS .; \
	fi

##@ Deployment

.PHONY: install
install: kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found -f -

.PHONY: deploy
deploy: kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found -f -

##@ Tools

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	@test -s $(LOCALBIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: clean
clean: ## Clean build artifacts.
	rm -rf bin/ cover.out
