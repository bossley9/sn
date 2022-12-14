PREFIX = /usr/local
BIN = $(PREFIX)/bin
EXE = sn

build: test
	go build -tags production -o ./$(EXE) ./main.go

test:
	go test ./pkg/jsondiff

run:
	go run ./main.go
run-c:
	go run ./main.go c
run-d:
	go run ./main.go d
run-h:
	go run ./main.go h
run-r:
	go run ./main.go r
run-u:
	go run ./main.go u

clean:
	rm ./$(EXE)

install:
	mkdir -p $(BIN)
	cp -f ./$(EXE) $(BIN)
	chmod 555 $(BIN)/$(EXE)

uninstall:
	rm -f $(BIN)/$(EXE)

.PHONY: all build clean install uninstall
