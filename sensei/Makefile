.PHONY: build test clean run

build:
	CGO_ENABLED=1 go build -o sensei ./cmd/sensei

test:
	go test ./internal/...

run: build
	./sensei

clean:
	rm -f sensei
