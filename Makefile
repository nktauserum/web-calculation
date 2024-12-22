build:
	@go build -o ./bin/main cmd/main.go
run: build
	@./bin/main

test: build
	@go test ./test/*