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

.PHONY: vendor
vendor: ### Vendor dependencies
	@go mod vendor

.PHONY: deps
deps:	### Optimize dependencies
	@go mod tidy

.PHONY: codegen
codegen: vendor ### Generate code
	@bash ./codegen/codegen.sh

.PHONY: fmt
fmt: ### Format
	@gofmt -s -w .

.PHONY: vet
vet: ### Vet
	@go vet ./...

### Lint
.PHONY: lint
lint: fmt vet

### Clean test 
.PHONY: test-clean
test-clean: ### Clean test cache
	@go clean -testcache ./...

.PHONY: test
test: lint ### Run tests
	@go test -v  -coverprofile=cover.out -timeout 10s ./...

.PHONY: cover
cover: test ### Run tests and generate coverage
	@go tool cover -html=cover.out -o=cover.html

.PHONY: mocks
mocks: ### Generate mocks
	@mockery --all --output internal/mocks

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