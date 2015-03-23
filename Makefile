default: test

deps:
	go get ./...

build: deps

test: build
	go test ./...
