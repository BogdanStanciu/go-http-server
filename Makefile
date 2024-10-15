.DEFAULT_GOAL := help

# Take all arguments from command-line and convert 
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

help:
	@fgrep -h "##" $(MAKEFILE_LIST)| sort | fgrep -v fgrep | tr -d '##' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##build: @ Build the project
build:
	@go build -o bin/go-http-server src/*.go  

##run: @ Run the project
run:
	@./bin/go-http-server $(RUN_ARGS)

##test: @ Run the tests
test: 
	@go test ./... -v


	