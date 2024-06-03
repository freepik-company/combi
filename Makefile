
.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
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

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/combi cmd/combi/main.go

RUN_SUBCOMMAND?=daemon

RUN_SYNC_FLAGS?=--sync-time=5s
RUN_SRC_FLAGS?=--source-path=config/samples/libconfig.yaml --source-type=git --source-field=example1
RUN_GIT_FLAGS?=--git-ssh-url=git@github.com:sebastocorp/combi.git --git-branch=main
RUN_FLAGS?=$(RUN_SYNC_FLAGS) $(RUN_SRC_FLAGS) $(RUN_GIT_FLAGS)

.PHONY: run
run: fmt vet ## Run a command from your host (need to be defined envs: RUNSUBCOMMAND and RUNFLAGS).
	go run cmd/combi/main.go $(RUN_SUBCOMMAND) $(RUN_FLAGS)
