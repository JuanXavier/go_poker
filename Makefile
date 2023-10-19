build:

run:
	@go build -o bin/go_poker
	@./bin/go_poker

test:
	go test -v ./...
