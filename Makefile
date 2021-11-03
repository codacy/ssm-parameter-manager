.PHONY: build
build:
	go build -o bin/ssm-parameter-manager

.PHONY: test
test:
	go test -v ./internal