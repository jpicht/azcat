MULTI=$(shell ls cmd/multi)
SOURCES=$(shell find pkg internal -name \*.go)

all: ${MULTI}

${MULTI}: cmd/multi/$@ $(shell find cmd/multi -name \*.go) ${SOURCES}
	go build ./cmd/multi/$@

azblob: $(shell find cmd/azblob -name \*.go) ${SOURCES}
	go build ./cmd/azblob
