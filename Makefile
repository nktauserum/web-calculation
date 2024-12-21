build:
	@go build -o ./bin/main cmd/main.go
run: build
	@./bin/main

test: build
	@./bin/main > /dev/null & ./test/test.sh