.PHONY: build build-sensei build-dojo test clean run run-dojo

build: build-sensei build-dojo

build-sensei:
	CGO_ENABLED=1 go build -o sensei ./cmd/sensei

build-dojo:
	go build -o dojo ./cmd/dojo

test:
	go test ./internal/...

run: build
	./sensei

run-dojo: build-dojo
	./dojo

clean:
	rm -f sensei dojo
