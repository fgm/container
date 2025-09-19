all: lint build

.PHONY: lint test build
lint:
	go tool staticcheck ./...

.PHONY: fuzz-smoke
fuzz-smoke: lint
	# Smoke tests on fuzzing targets.
	go test -fuzz='\QFuzzBasicMapAdd\E'   -fuzztime=10s ./set
	go test -fuzz='\QFuzzBasicMapItems\E' -fuzztime=10s ./set
	go test -fuzz='\QFuzzBasicMapUnion\E' -fuzztime=10s ./set

.PHONY: test
test:
	# This includes the fuzz tests in unit test mode
	go test -race ./...

.PHONY: build
build: test fuzz-smoke
	go build ./...
