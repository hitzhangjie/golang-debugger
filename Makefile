.PHONY: all

all:
	GOPATH=$(shell go env GOPATH):$(shell pwd) GOOS=linux GOARCH=amd64 go build -o debugger src/main.go

