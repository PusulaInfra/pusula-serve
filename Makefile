.PHONY: all tidy build test run clean docker-up

# Varsayılan hedef: Paketleri düzenle ve projeyi derle
all: tidy build test

# 1. Go modüllerini düzenle ve eksik bağımlılıkları indir
tidy:
	@echo "==> Go modülleri senkronize ediliyor..."
	go mod tidy

# 2. Projeyi derle (Build)
build: tidy
	@echo "==> PusulaInfra motoru derleniyor..."
	go build -v -o pusula-serve ./cmd/pusula-serve

# 3. Testleri ve motor entegrasyonunu doğrula
test:
	@echo "==> Testler çalıştırılıyor..."
	go test -v ./engine/...

# 4. CLI Modunda Hızlı Analiz Çalıştır (Örn: Llama-3.1 70B & 4x H100)
run-cli: build
	@echo "==> CLI Analiz Modu Başlatılıyor..."
	./pusula-serve -cli -model llama-3.1-70b -gpus 4 -provider lambda

# 5. HTTP Servis ve Ops Sunucusunu Başlat
run-server: build
	@echo "==> Enterprise HTTP Sunucusu :8080 portunda ayağa kaldırılıyor..."
	./pusula-serve -addr :8080

# 6. Docker Compose ile Tüm Sistemi (Prometheus dahil) Tek Komutla Uçur
docker-up:
	@echo "==> Docker Compose stack ayağa kaldırılıyor..."
	docker compose up --build -d

# 7. Derleme artıklarını temizle
clean:
	@echo "==> Temizlik yapılıyor..."
	rm -f pusula-serve
	go clean
