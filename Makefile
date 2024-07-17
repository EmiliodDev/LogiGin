build:
	@go build -o bin/todoAPI cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/todoAPI 
