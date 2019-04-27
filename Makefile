APPLICATIONS=$(shell find cmd/ -type f -name main.go | xargs -I file echo "file" | rev | cut -d'/' -f2 | rev)
TESTS=$(shell find . -type d -name test)
GO_OS=linux
ifeq ($(OS),Windows_NT)
	GO_OS=windows
endif

.PHONY: all deps format test build ci-setup

all: deps test build

deps:
	@echo "Getting Deps"; \
	GOOS=${GO_OS} go get "github.com/google/gopacket"; \
	GOOS=${GO_OS} go get "github.com/google/gopacket/pcap"; \
	GOOS=${GO_OS} go get "github.com/artman41/vJoy"; \
	GOOS=${GO_OS} go get "github.com/gizak/termui"; \
	GOOS=${GO_OS} go get "github.com/gizak/termui/widgets";

format:
	@echo "Formatting Code"; \
	for file in $(echo "`find . -name '*.go' | xargs -I file dirname 'file'`" | sort | uniq); do \
		go vet "$$file"; \
		go fmt "$$file"; \
	done;

test: format
	@echo "Testing Packages"; \
	for dir in ${TESTS}; do \
		go test $$dir; \
	done;

build: format
	@echo "Building Binaries"; \
	mkdir -p cmd/bin; \
	for app in ${APPLICATIONS}; do \
		APP_NAME=`grep "%name%" "cmd/$$app/.info" | awk '{print $$2}'`
		APP_VERSION=`grep "%version%" "cmd/$$app/.info" | awk '{print $$2}'`
		GOOS=windows GOARCH=386 go build "cmd/$$app/*.go" -o "cmd/bin/$$APP_NAME-$$APP_VERSION_x86.exe"; \
		GOOS=windows GOARCH=amd64 go build "cmd/$$app/*.go" -o "cmd/bin/$$APP_NAME-$$APP_VERSION_x64.exe"; \
	done;

ci-setup:
	@printenv && echo; \
	echo "`git config --get remote.origin.url`"; \
	AuthorRepo=`git config --get remote.origin.url | sed "s|https://github.com/||g" | cut -d ':' -f2 | rev | cut -d '.' -f2 | rev | sed "s|com/||g"`; \
	echo "AuthorRepo: $$AuthorRepo"; \
	Author=`echo "$$AuthorRepo" | cut -d '/' -f1`; \
	echo "Author: $$Author"; \
	Repo=`echo "$$AuthorRepo" | cut -d '/' -f2`; \
	echo "Repo: $$Repo"; \
	mkdir -p "$$GOPATH/src/github.com/$$Author"; \
	cd ..; \
	mv $$Repo "$$GOPATH/src/github.com/$$Author/";