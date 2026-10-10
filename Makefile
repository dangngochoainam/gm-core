# A fresh clone needs Go and nothing else. buf and the three code generators are
# pinned as tool dependencies in tools/go.mod, checksummed in tools/go.sum, and
# built into bin/ by the rules below; the first `make gen-proto` on a machine does
# that bootstrap on its own, and later runs skip it. Nothing is installed onto the
# host, and generating needs no network once bin/ exists.
#
# tools/ is a module of its own because those four tools drag in about a hundred
# dependencies. Declared in gm-core/go.mod they would land in the module graph of
# every service, since every service requires gm-core.
#
# The generators run as local plugins rather than remote ones. Remote plugins call
# the Buf Schema Registry on every single run, not just the first, so they make
# codegen depend on the network being up. Same pinned versions either way: the
# output is byte-identical.
BIN := bin

PKG_buf                     := github.com/bufbuild/buf/cmd/buf
PKG_protoc-gen-go           := google.golang.org/protobuf/cmd/protoc-gen-go
PKG_protoc-gen-go-grpc      := google.golang.org/grpc/cmd/protoc-gen-go-grpc
PKG_protoc-gen-grpc-gateway := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

TOOLS := $(BIN)/buf $(BIN)/protoc-gen-go $(BIN)/protoc-gen-go-grpc $(BIN)/protoc-gen-grpc-gateway

# Stem of the target is the binary name; PKG_<name> maps it to the package to
# build. Rebuilds when the pinned versions change, because both go.mod and go.sum
# are prerequisites.
$(BIN)/%: tools/go.mod tools/go.sum
	@mkdir -p $(BIN)
	cd tools && go build -o ../$@ $(PKG_$*)

tools: $(TOOLS)

# The input is `proto`, not `.`, so third_party/ stays an import root and never
# becomes a generation target. buf.yaml is what makes both visible as one workspace.
gen-proto: $(TOOLS)
	./$(BIN)/buf generate proto

mod:
	go mod tidy

build:
	go build ./...

# golangci-lint pin theo version, chạy qua `go run` nên không cài gì lên máy.
# Không nằm trong tools/ vì nó cần go >= 1.26, còn mọi module giữ go 1.24.0.
# Config ở .golangci.yml; chỉ lỗi mới so với origin/master mới làm fail.
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0

fmt:
	$(GOLANGCI_LINT) fmt ./...

lint:
	$(GOLANGCI_LINT) run ./...

test:
	go test ./...

verify: lint build

.PHONY: tools gen-proto mod build fmt lint test verify
