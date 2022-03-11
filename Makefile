ALL=$(shell ls cmd)
SOURCES=$(shell find pkg internal -name \*.go)

all: ${ALL}

${ALL}: cmd/$@ ${SOURCES}
	go build ./cmd/$@
