build:
	@go build -o bin/go-http-server src/*.go

run: build
	@./bin/go-http-server


test: 
	@go test ./... -v


	