
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

# Sync time flags
DAEMON_SYNC_TIME?=5s
DAEMON_SYNC_FLAGS?=--sync-time=$(DAEMON_SYNC_TIME)
# Source config flags
DAEMON_SRC_TYPE?=local
DAEMON_SRC_PATH?=config/samples/libconfig.yaml
DAEMON_SRC_FIELD?=example1
DAEMON_SRC_FLAGS?=--source-type=$(DAEMON_SRC_TYPE) --source-path=$(DAEMON_SRC_PATH) --source-field=$(DAEMON_SRC_FIELD)
# Extra config flags
DAEMON_EXTRA_FLAGS?=
# Daemon subcommand flags
DAEMON_FLAGS?=$(DAEMON_SYNC_FLAGS) $(DAEMON_SRC_FLAGS) $(DAEMON_EXTRA_FLAGS)

.PHONY: run-daemon
run-daemon: fmt vet ## Run a command from your host (define DAEMON_FLAGS to custom run daemon).
	go run cmd/combi/main.go daemon $(DAEMON_FLAGS)
