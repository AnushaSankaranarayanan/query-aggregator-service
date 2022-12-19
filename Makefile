.PHONY: all test build clean cover

all: clean test build

build: 
	mkdir -p build
	go build -o build -tags real ./...

test:
	go test -v -coverprofile=tests/results/cover.out -tags fake ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...


