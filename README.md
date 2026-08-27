# Pusula Serve 🚀
*LLM Serving Config & Optimization Copilot*

**Pusula Serve**, yapay zeka mühendislerinin büyük dil modellerini (LLM) cluster üzerinde en yüksek verimlilikle sunabilmeleri (`vLLM` & `SGLang`) için geliştirilmiş açık kaynaklı bir Go aracıdır.

## ✨ Öne Çıkan Özellikler
- **Multi-Engine Desteği:** İster **vLLM** (High-Throughput) ister **SGLang** (RadixAttention & Low-Latency) motoru seçerek optimize komutlar üretin.
- **Akıllı Bellek & VRAM Analizi:** Model parametreleri ve context uzunluğuna göre anlık VRAM tüketimi tahmini ve cluster tavsiyesi.
- **İnteraktif Web Paneli:** Modern koyu tema tasarımıyla tarayıcı üzerinden kolay konfigürasyon yönetimi.

## 🛠️ Modüler Mimari
- `main.go`: Sunucu tetikleyici ana giriş noktası.
- `parser.go`: Bellek, VRAM ve maliyet hesaplama optimizasyon motoru.
- `server.go`: İnteraktif HTTP web sunucusu ve arayüz katmanı.

## 🚀 Çalıştırma
Projeyi yerel ortamınızda ayağa kaldırmak için:
```bash
go run main.go parser.go server.go
