# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

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

# .PHONY: test
# test: fmt vet ## Run tests.

##@ Build

.PHONY: build
build: fmt vet ## Build binary.
	go build -o bin/kubectl-lethe main.go

##@ Deployment

.PHONY: install
install: build ## binary to your GOBIN.
	cp bin/kubectl-lethe $(GOBIN)

.PHONY: uninstall
uninstall: ## remove binary from your GOBIN
	rm $(GOBIN)/kubectl-lethe

# CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
# .PHONY: controller-gen
# controller-gen: ## Download controller-gen locally if necessary.
# 	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.2)

# KUSTOMIZE = $(shell pwd)/bin/kustomize
# .PHONY: kustomize
# kustomize: ## Download kustomize locally if necessary.
# 	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# ENVTEST = $(shell pwd)/bin/setup-envtest
# .PHONY: envtest
# envtest: ## Download envtest-setup locally if necessary.
# 	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

# # go-get-tool will 'go get' any package $2 and install it to $1.
# PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
# define go-get-tool
# @[ -f $(1) ] || { \
# set -e ;\
# TMP_DIR=$$(mktemp -d) ;\
# cd $$TMP_DIR ;\
# go mod init tmp ;\
# echo "Downloading $(2)" ;\
# GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
# rm -rf $$TMP_DIR ;\
# }
# endef
