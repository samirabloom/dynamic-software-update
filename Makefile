# Makefile for a go project
#
# Author: Samira Rabbanian
# 	
# Targets:
# 	all:          cleans and builds all the code
# 	clean:        cleans the code
#	test:         runs all the tests
#	coverage:     runs all the tests and outputs a coverage report
# 	build:        builds all the code
# 	dependencies: installs dependent projects
# 	install:      installs the code to /usr/local/bin
# 	run:          the proxy with INFO log level

GOBIN := /usr/local/bin/
GOPATH := $(shell pwd):$(GOPATH)
PATH := $(PATH):$(GOPATH)/bin

FLAGS := GOPATH=$(GOPATH)


all: install

clean:
	$(FLAGS) go clean -i -x ./.../$*
	rm -rf $(GOBIN)proxy proxy pkg

test: clean
	vagrant up docker
	$(FLAGS) go test -v ./.../$*
	
coverage:
	go get github.com/axw/gocov/gocov
	go get gopkg.in/matm/v1/gocov-html
	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy | gocov-html > proxy_coverage.html
#	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy/contexts | gocov-html > contexts_coverage.html
#	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy/log | gocov-html > log_coverage.html
	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy/stages | gocov-html > stages_coverage.html
	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy/tcp | gocov-html > tcp_coverage.html
#	PATH=$(PATH):$(GOPATH)/bin gocov test -v proxy/transition | gocov-html > transition_coverage.html

dependencies:
	go get -v code.google.com/p/go-uuid/uuid
	go get -v github.com/op/go-logging
	go get -v github.com/fsouza/go-dockerclient
	go get -v github.com/franela/goreq

build: clean dependencies
	$(FLAGS) go build -v -o proxy ./src/main_run.go

count:
	find . -name "*.go" -print0 | xargs -0 wc -l

install: build
	cp proxy $(GOBIN)

run: install
	proxy -logLevel INFO
