EXE = sn

build:
	go build -o ./$(EXE) ./cmd/main.go

run:
	go run ./cmd/main.go
