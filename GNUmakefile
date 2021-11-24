BUILD_ID = $(shell git describe --tags HEAD)

BIN = evcli

LDFLAGS = -X main.buildId=$(BUILD_ID)

all: build

build: FORCE
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN)

test:
	go test -race ./...

vet:
	go vet ./...

clean:
	$(RM) $(BIN)

FORCE:

.PHONY: all build test vet clean
