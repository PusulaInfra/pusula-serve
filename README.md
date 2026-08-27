# Pusula Serve 🚀
*Enterprise LLM Serving Config & Optimization Copilot*

**Pusula Serve**, yapay zeka mühendislerinin büyük dil modellerini (LLM) cluster üzerinde en yüksek verimlilikle sunabilmeleri (`vLLM` & `SGLang`), otomatik donanım dağılımı yapabilmeleri ve Kubernetes YAML manifestoları üretebilmeleri için Go ile geliştirilmiş kurumsal bir açık kaynak aracıdır.

## ✨ Kurumsal Öne Çıkan Özellikler
- **Hazır Model Presetleri:** Llama-3-70B, DeepSeek-V3 (MoE), Qwen-2.5-72B ve Mistral-Large gibi dev modelleri tek tıkla seçebilme.
- **Multi-Engine Desteği:** İster **vLLM** (High-Throughput) ister **SGLang** (RadixAttention & Low-Latency) motoruyla optimize komutlar üretme.
- **Akıllı TP & PP Dağılım Matrisi:** GPU sayısına ve context uzunluğuna göre Tensor Parallelism (TP) ve Pipeline Parallelism (PP) değerlerini otomatik hesaplama.
- **Kubernetes (K8s) YAML Üretici:** Üretilen konfigürasyonu doğrudan kurumsal k8s deployment manifestosuna dönüştürme.
- **Çok Dilli Web Paneli (i18n):** Türkçe ve İngilizce anlık dil desteği.

## 🛠️ Modüler Mimari
- `main.go`: Sunucu tetikleyici ana giriş noktası.
- `parser.go`: Bellek, VRAM ve maliyet hesaplama optimizasyon motoru.
- `server.go`: İnteraktif HTTP web sunucusu, i18n ve kurumsal arayüz katmanı.

## 🚀 Çalıştırma
Projeyi yerel ortamınızda ayağa kaldırmak için:
```bash
go run main.go parser.go server.go
