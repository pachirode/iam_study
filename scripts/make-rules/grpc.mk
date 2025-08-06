PROTOC_INC_PATH=$(dir $(shell which protoc 2>/dev/null))/../include
API_DEPS=$(ROOT_DIR)/internal/pkg/api/proto/apiserver/v1/cache.proto
GAPI_ROOT=$(ROOT_DIR)/internal/pkg/api/proto/apiserver/v1
API_DEPSRCS=$(API_DEPS:.proto=.pb.go)

.PHONY: gen
gen: gen.clean gen.protoc

.PHONY: gen.protoc
gen.protoc: gen.plug.verify
	@echo "================> Generate protobuf files"
	@protoc \
		--proto_path=$(GAPI_ROOT) \
		--go_out=paths=source_relative:$(GAPI_ROOT) \
		--go-grpc_out=paths=source_relative:$(GAPI_ROOT) \
		$(shell find $(GAPI_ROOT) -name *.proto)

.PHONY: gen.plug.verify
gen.plug.verify:
ifeq (,$(shell which protoc-gen-go 2>/dev/null))
	@echo "=================> Installing protoc-gen-go"
	@GO111MODULE=on $(GO) get install github.com/golang/protobuf/protoc-gen-go@latest
endif

.PHONY: gen.clean
gen.clean:
	@rm -f $(API_DEPSRCS)
