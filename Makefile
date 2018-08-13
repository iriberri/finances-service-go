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
	go clean -cache

check: deps clean test

ci-trigger:
	#
	# Ensuring dependencies are installed
	#
	dep ensure -v
	#
	# Tests
	#
	go test ./src/finances-service/... -v
