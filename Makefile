SHELL = /bin/bash

PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Setting GOBIN and PATH ensures two things:
# - All 'go install' commands we run
#   only affect the current directory.
# - All installed tools are available on PATH
#   for commands like go generate.
export GOBIN = $(PROJECT_ROOT)/bin
export PATH := $(GOBIN):$(PATH)

TEST_FLAGS ?= -v -race

STRINGER = bin/stringer
TOOLS = $(STRINGER)
