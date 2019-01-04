
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


build:

	go build -o 181231-web main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t 181231-web:latest
