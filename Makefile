MULTI=$(shell ls cmd/multi)
SOURCES=$(shell find auth actions internal -name \*.go)

all: ${MULTI} azblob

${MULTI}: cmd/multi/$@ $(shell find cmd/multi -name \*.go) ${SOURCES}
	go build ./cmd/multi/$@

azblob: $(shell find cmd/azblob -name \*.go) ${SOURCES}
	go build ./cmd/azblob

clean:
	rm -f ${MULTI} azblob
