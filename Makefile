.PHONY: none
none:


deps: go.mod
	@echo "+ $@"
	@go mod tidy
ifdef CI
	@git diff --exit-code -- go.mod go.sum || { echo "go.mod/go.sum files were updated after running 'go mod tidy', run this command on your local machine and commit the results." ; exit 1 ; }
endif
	go mod verify
	@touch deps

GOBIN := $(CURDIR)/.gobin
PATH := $(GOBIN):$(PATH)

# Makefile on Mac doesn't pass the updated PATH and GOBIN to the shell
# and so, without the following line, the shell does not end up
# trying commands in $(GOBIN) first.
# See https://stackoverflow.com/a/36226784/3690207
SHELL := env GOBIN="$(GOBIN)" PATH="$(PATH)" /bin/bash

########################################
###### Binaries we depend on ###########
########################################

GOLANGCILINT_BIN := $(GOBIN)/golangci-lint
$(GOLANGCILINT_BIN): deps
	@echo "+ $@"
	cd tool-imports; \
	GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint

###########
## Lint ##
###########

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT_BIN)
ifdef CI
	@echo '+ $@'
	@echo 'The environment indicates we are in CI; running linters in check mode.'
	@echo 'If this fails, run `make lint`.'
	golangci-lint run
else
	golangci-lint run --fix
endif

.PHONY: lint
lint: golangci-lint

####################
## Code generation #
####################

.PHONY: go-generated-srcs
go-generated-srcs: deps
	go generate ./...

.PHONY: generated-srcs
generated-srcs: go-generated-srcs

#############
## Compile ##
#############


.PHONY: build
build:

##########
## Test ##
##########

.PHONY: test
test:
	go test ./...
