# "go install"-ed binaries will be placed here during development.
export GOBIN ?= $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

GO_FILES = $(shell find . \
	   -path '*/.*' -prune -o \
	   '(' -type f -a -name '*.go' ')' -print)

GOLINT = bin/golint
STATICCHECK = bin/staticcheck

TOOLS = $(GOLINT) $(STATICCHECK)

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: tools
tools: $(TOOLS)

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: gofmt golint staticcheck

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || \
		(echo "gofmt failed. Please reformat the following files:" | \
		cat - $(FMT_LOG) && false)

.PHONY: golint
golint: $(GOLINT)
	golint ./...

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	staticcheck ./...

$(GOLINT): tools/go.mod
	cd tools && go install golang.org/x/lint/golint

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck
