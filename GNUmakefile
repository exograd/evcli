BIN = evcli

all: build

build: FORCE
	go build -o $(BIN)

test:
	go test -race ./..

clean:
	$(RM) $(BIN)

FORCE:

.PHONY: all build test clean
