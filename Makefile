VERSION := $(shell cat .version)
LDFLAGS := -linkmode external -w -extldflags "-static" -X main.version=${VERSION}

.PHONY: build
build:
	go build -o bin/ssm-parameter-manager  -ldflags '$(LDFLAGS)'

.PHONY: test
test:
	go test -v ./sops ./ssm