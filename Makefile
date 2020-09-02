GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG := $(shell git describe --abbrev=0 HEAD 2>/dev/null)
LD_FLAGS := '-s -w \
	-X main.gitTag=$(GIT_TAG) \
	-X main.gitRef=$(GIT_SHA) \
	-X main.shellCompletion=$(shell base64 -w0 watchmon-completion)'

$(GOPATH)/bin/watchmon: bin/watchmon
	cp bin/watchmon $(GOPATH)/bin

bin/watchmon: *.go deps.txt
	go build \
		-v \
		-o $@ \
		-ldflags=$(LD_FLAGS) \
		-trimpath

deps.txt: go.mod go.sum
	go get
	go mod graph > deps.txt
