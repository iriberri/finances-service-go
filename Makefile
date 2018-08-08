default: test build

deps:
	dep ensure

build:
	go build ./src/finances-service

test:
	go test ./src/finances-service/...

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
