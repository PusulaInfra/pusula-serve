<div align="center">
  <img src="https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=1200&auto=format&fit=crop" alt="Pusula Serve Banner" width="100%" style="border-radius: 12px;"/>
  
  # Pusula Serve 🚀
  *Enterprise LLM Serving Config & Optimization Copilot*

  <p><b>Go</b> ile geliştirilmiş; vLLM & SGLang destekli, bulut maliyet analizi ve Kubernetes manifestosu üretebilen kurumsal LLM Ops SaaS platformu.</p>

  <p>
    <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Engine-vLLM%20%7C%20SGLang-orange?style=flat-square" alt="Engines">
    <img src="https://img.shields.io/badge/Edition-Elite%20SaaS-blueviolet?style=flat-square" alt="Edition">
    <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  </p>
</div>

---

## ✨ Kurumsal Öne Çıkan Özellikler
- **Hazır Model Presetleri:** Llama-3-70B, DeepSeek-V3 (MoE), Qwen-2.5-72B ve Mistral-Large gibi dev modelleri tek tıkla seçebilme.
- **Multi-Engine Desteği:** İster **vLLM** (High-Throughput) ister **SGLang** (RadixAttention & Low-Latency) motoruyla optimize komutlar üretme.
- **Akıllı TP & PP Dağılım Matrisi:** GPU sayısına ve context uzunluğuna göre Tensor Parallelism (TP) ve Pipeline Parallelism (PP) değerlerini otomatik hesaplama.
- **Bulut Altyapı Maliyet Simülasyonu:** Çoklu GPU cluster ekonomisini anlık hesaplama ve maliyet optimizasyonu sunma.
- **Kubernetes (K8s) YAML Üretici:** Üretilen konfigürasyonu doğrudan kurumsal k8s deployment manifestosuna dönüştürme.
- **Çok Dilli Web Paneli (i18n):** Türkçe ve İngilizce anlık dil desteği.

## 🛠️ Modüler Mimari
- `main.go`: Sunucu tetikleyici ana giriş noktası.
- `parser.go`: Bellek, VRAM ve maliyet hesaplama optimizasyon motoru.
- `server.go`: İnteraktif HTTP web sunucusu, i18n, bulut simülasyonu ve kurumsal arayüz katmanı.

## 🚀 Çalıştırma
Projeyi yerel ortamınızda ayağa kaldırmak için:
```bash
go run main.go parser.go server.go
