SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif

ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(ROOT_DIR)/tools
endif

ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
endif

ifeq ($(origin VERSION),undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

GIT_TREE_STATE := "dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif

GIT_COMMIT := $(shell git rev-parse HEAD)

ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

PLATFORMS ?= linux_amd64 linux_arm64

ifeq ($(origin PLATFORM),undefined)
	ifeq ($(origin GOOS),undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH),undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	IMAGE_PLAT := linux_$(GOARCH)
else
		GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
		GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
		IMAGE_PLAT := $(PLATFORM)
endif

FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty

ifndef v
MAKEFLAGS += --no-print-directory
endif

ifeq ($(origin CERTIFICATES),undefined)
CERTIFICATES=iam-apiserver
endif

BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint  goimports codegen
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release go-mod-outdated protoc-gen-go cfssl go-gitlint
TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score

COMMA := ,
SPACE :=
SPACE +=
