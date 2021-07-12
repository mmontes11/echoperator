UNAME := $(uname -s)
LD_FLAGS := -X main.version=$(VERSION) -s -w

export CGO_ENABLED=0

.PHONY: all
all: help

.PHONY: help
help:	### Show targets documentation
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: deps
deps:	### Get dependencies
	@go mod tidy

.PHONY: vendor
vendor: ### Vendor dependencies
	@go mod vendor

.PHONY: codegen
codegen: vendor ### Generate code
	@bash ./codegen/codegen.sh

.PHONY: run  
run: ### Run example
	@go run cmd/main.go

.PHONY: clean
clean: ### Clean build files
	@rm -rf ./bin
	@go clean

.PHONY: build
build: clean ### Build binary
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ./bin/echoperator ./cmd/*.go
	@chmod +x ./bin/*