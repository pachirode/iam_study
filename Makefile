.DEFAULT_GOAL := all

.PHONY: all
all: tidy format

# ==========================================================================================
# Build options

ROOT_PACKAGE=github.com/pachirode/iam_study
VERSION_PACKAGE=github.com/pachirode/iam_study/pkg/version

# ==========================================================================================
# Includes

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/swagger.mk
include scripts/make-rules/ca.mk
include scripts/make-rules/grpc.mk

.PHONY: clean
clean:
	@echo "============> Clean all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "============> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

.PHONY: tools
tools:
	@$(MAKE) tools.install

.PHONY: tidy
tidy:
	@$(GO) mod tidy
