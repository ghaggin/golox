all: run

SRC := $(shell find . -type f -name "*.go" ! -name "*_test.go")

.PHONY: test
run:
	@go run ${SRC}

test:
	go test -v *.go

file:
	@go run ${SRC} test.lox
