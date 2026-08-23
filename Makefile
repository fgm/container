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

.PHONY: bench
bench:
	# Run this when modifying the code to obtain data to update BENCHMARKS.md
	go test -bench=Benchmark ./...

.PHONY: cover
cover:
	# This runs the benchmarks just once, as unit tests, for coverage reporting only.
	# It does not replace running "make bench".
	mkdir -p coverage
	go test -v -race -run=. -bench=. -benchtime=1x -coverprofile=coverage/cover.out -covermode=atomic ./...

.PHONY: test
test:
	# This includes the fuzz tests in unit test mode
	go test -race ./...

.PHONY: build
build: test fuzz-smoke
	go build ./...
