default: test build

deps:
	dep ensure

build:
	go build ./src/finances-service

test:
	go test ./src/finances-service/...

run:
	go run src/finances-service/main.go

clean:
	#
	# Cleaning the go environment for this project,
	# cache and test cache only.
	#
	go clean -cache -testcache ./...

clean-all:
	#
	# Cleaning the go environment for this project,
	# as radically as possible.
	#
	go clean -cache -testcache -x -r -i ./...

check: deps clean test

ci-trigger: clean-all
	#
	# Ensuring dependencies are installed
	#
	dep ensure -v
	#
	# Build the project
	#
	go build ./src/finances-service
	#
	# Vet
	#
	go vet -v ./src/finances-service/...
	#
	# Tests
	#
	go test -v ./src/finances-service/...
