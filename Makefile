BIN := $(shell basename $$PWD).bin

LAST_COMMIT := $(shell git rev-parse --short HEAD)
LAST_COMMIT_DATE := $(shell git show -s --format=%ci ${LAST_COMMIT})
BUILD_DATE := ${VERSION} (${LAST_COMMIT} $(shell date -u +"%Y-%m-%dT%H:%M:%S%z"))
VERSION := $(shell git describe)

build:
	@go build -a -o ${BIN} -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'"
	$(info Build successful. Current build version: $(VERSION))
.PHONY: build

staticcheck:
	@go install honnef.co/go/tools/cmd/staticcheck
	@ $(GOPATH)/bin/staticcheck -tests=false ./...

test:
	@go test ./...
.PHONY: test

run: build
	@./${BIN}
.PHONY: run

clean:
	@go clean
	-@rm -f ${BIN}
.PHONY: clean
