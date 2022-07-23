EXE = sn

build: test
	go build -o ./$(EXE) ./cmd/main.go

test:
	go test ./pkg/jsondiff ./pkg/client

run:
	go run ./cmd/main.go
