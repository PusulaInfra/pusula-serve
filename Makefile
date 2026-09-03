# Pusula Serve — CLI, API, and Pages console are one product.
.PHONY: all tidy build test run run-cli run-server clean docker-up

all: tidy build test

tidy:
	go mod tidy

build: tidy
	go build -v -o pusula-serve ./cmd/pusula-serve

test:
	go test ./...

run: run-cli

run-cli: build
	./pusula-serve -cli -model llama-3.3-70b -gpu H100 -gpus 4 -ctx 16384 -seqs 16 -provider lambda

run-server: build
	./pusula-serve -addr :8080

docker-up:
	docker compose up --build -d

clean:
	rm -f pusula-serve
	go clean
