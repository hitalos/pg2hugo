.PHONY: build

build:
	go build -ldflags '-s -w' ./cmd/pg2hugo
