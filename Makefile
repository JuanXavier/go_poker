build:
	@go build -o bin/go_poker

run:
	@./bin/go_poker

test:
	go test -v ./...
