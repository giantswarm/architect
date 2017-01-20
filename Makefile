all: architect
.PHONY: image test

architect: main.go
	docker run \
		-v $(shell pwd):/go/src/github.com/giantswarm/architect \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOPATH=/go \
		-e CGOENABLED=0 \
		-w /go/src/github.com/giantswarm/architect \
		golang:1.7.4 \
		go build -a -v -tags netgo

image: architect Dockerfile
	docker build --rm=false -t registry.giantswarm.io/architect:test .
