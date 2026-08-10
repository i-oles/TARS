run:
	go run cmd/main.go

lint:
	golangci-lint run ./...

test:
	go test -v ./...

build:
	go build -v -o bin/tars cmd/main.go

