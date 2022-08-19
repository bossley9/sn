PREFIX = /usr/local
BIN = $(PREFIX)/bin
EXE = sn

build: test
	go build -tags isproduction -o ./$(EXE) ./cmd/main.go

test:
	go test ./pkg/jsondiff

run:
	go run ./cmd/main.go

run-c:
	go run ./cmd/main.go c

run-r:
	go run ./cmd/main.go r

run-u:
	go run ./cmd/main.go u

clean:
	rm ./$(EXE)

install:
	mkdir -p $(BIN)
	cp -f ./$(EXE) $(BIN)
	chmod 555 $(BIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)

.PHONY: all build clean install uninstall
