BINARY_NAME := ssm-parameter-manager
VERSION := $(shell cat .version)
LDFLAGS_STATIC := -linkmode external -w -extldflags "-static" 
LDFLAGS := -X main.version=${VERSION}

.PHONY: build
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}_linux_amd64 -ldflags '$(LDFLAGS_STATIC) ${LDFLAGS}'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}_darwin_amd64 -ldflags '${LDFLAGS}'

.PHONY: test
test:
	go test -v ./sops ./ssm