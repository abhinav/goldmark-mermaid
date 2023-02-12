# "go install"-ed binaries will be placed here during development.
export GOBIN ?= $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

GO_FILES = $(shell find . \
	   -path '*/.*' -prune -o \
	   '(' -type f -a -name '*.go' ')' -print)

GOLINT = bin/golint
STATICCHECK = bin/staticcheck
STRINGER = bin/stringer

TOOLS = $(GOLINT) $(STATICCHECK) $(STRINGER)

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

.PHONY: generate
generate: tools
	go generate -x ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: gofmt golint staticcheck check-generate

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
	$(STATICCHECK) ./...
	$(eval STATICCHECK_LOG := $(shell mktemp -t staticcheck.XXXXX))
	$(STATICCHECK) ./... | grep -v SA1019 > $(STATICCHECK_LOG) || true
	@[ ! -s "$(STATICCHECK_LOG)" ] || \
		(echo "static failed:" | \
		cat - $(STATICCHECK_LOG) && false)

.PHONY: check-generate
check-generate: generate
	@DIFF=$$(git diff --name-only); \
	if [ -n "$$DIFF" ]; then \
		echo "--- The following files are dirty:"; \
		echo "$$DIFF"; \
		exit 1; \
	fi

$(GOLINT): tools/go.mod
	cd tools && go install golang.org/x/lint/golint

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck

$(STRINGER): tools/go.mod
	cd tools && go install golang.org/x/tools/cmd/stringer
