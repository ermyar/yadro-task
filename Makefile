
all:
	go build ./cmd/biathlon/...

build:
	go build ./cmd/biathlon

test:
	go test -v ./...

clean:
	rm biathlon

.PHONY: test build
