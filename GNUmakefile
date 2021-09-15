BUILD_ID = $(shell git describe --tags HEAD)

BIN = evcli

LDFLAGS = -X main.buildId=$(BUILD_ID)

all: build

build: FORCE
	go build -ldflags "$(LDFLAGS)" -o $(BIN)

test:
	go test -race ./..

clean:
	$(RM) $(BIN)

FORCE:

.PHONY: all build test clean
