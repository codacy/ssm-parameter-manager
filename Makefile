VERSION := $(shell cat .version)
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: build
build:
	go build -o bin/ssm-parameter-manager  -ldflags $(LDFLAGS)

.PHONY: test
test:
	go test -v ./sops ./ssm