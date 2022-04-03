.PHONY: all
all: test bench

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: bench
bench:
	go test -bench=.