.PHONY: build run test bench clean docker pgo-build

# Gerçek üretim profiliyle (default.pgo) en üst düzey optimize derleme
build:
	go build -pgo=auto -ldflags="-s -w" -o bin/pusula-serve cmd/main.go

# Performans verisini toplamak için benchmark çalıştırma
profile:
	go test -bench=TestEngineStressAndQueue -cpuprofile=cpu.prof ./engine
	mv cpu.prof default.pgo

run:
	go run cmd/main.go

test:
	go test -v -race ./...

bench:
	go test -bench=. -benchmem ./engine

docker:
	docker build -t pusula-serve:latest .

clean:
	rm -rf bin/ default.pgo
