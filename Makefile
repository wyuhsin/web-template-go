GOHOSTOS:=$(shell go env GOHOSTOS) GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
API_PROTO_DIR := ./api

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
		--proto_path=./third_party \
		--go_out=paths=source_relative:./internal \
		$(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=$(API_PROTO_DIR) \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:$(API_PROTO_DIR) \
	       --go-http_out=paths=source_relative:$(API_PROTO_DIR) \
	       --go-grpc_out=paths=source_relative:$(API_PROTO_DIR) \
	       --go-errors_out=paths=source_relative:$(API_PROTO_DIR) \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       --validate_out=paths=source_relative,lang=go:$(API_PROTO_DIR) \
	       $(API_PROTO_FILES)

	@find $(API_PROTO_DIR) -name '*.pb.go' -type f | while read file; do \
		protoc-go-inject-tag -input=$$file; \
	done


.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;
