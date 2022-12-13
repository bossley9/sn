PREFIX = /usr/local
BIN = $(PREFIX)/bin
EXE = sn

build: test
	go build -tags production -o ./$(EXE) ./cmd/sn/main.go

test:
	go test ./pkg/jsondiff

run:
	go run ./cmd/sn/main.go
run-c:
	go run ./cmd/sn/main.go c
run-d:
	go run ./cmd/sn/main.go d
run-h:
	go run ./cmd/sn/main.go h
run-r:
	go run ./cmd/sn/main.go r
run-u:
	go run ./cmd/sn/main.go u

clean:
	rm ./$(EXE)

install:
	mkdir -p $(BIN)
	cp -f ./$(EXE) $(BIN)
	chmod 555 $(BIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)

.PHONY: all build clean install uninstall
