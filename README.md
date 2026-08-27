<div align="center">
  <img src="https://images.unsplash.com/photo-1558494949-ef010cbdcc31?q=80&w=1200&auto=format&fit=crop" alt="Pusula Serve Enterprise LLM Infrastructure" width="100%" style="border-radius: 12px;"/>
  
  # Pusula Serve 🚀
  *Enterprise LLM Serving Config & Optimization Copilot*

  <p><b>[EN]</b> An enterprise-grade open-source LLM Ops platform built in Go for high-throughput vLLM & SGLang deployment, cluster topology optimization, cloud cost simulation, and Kubernetes manifest generation.</p>
  <p><b>[TR]</b> vLLM & SGLang tabanlı, cluster optimizasyonu, bulut maliyet simülasyonu ve Kubernetes manifestosu üretebilen Go ile geliştirilmiş kurumsal LLM Ops aracı.</p>

  <p>
    <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Engines-vLLM%20%7C%20SGLang-orange?style=flat-square" alt="Engines">
    <img src="https://img.shields.io/badge/Edition-Elite%20SaaS-blueviolet?style=flat-square" alt="Edition">
    <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  </p>
</div>

---

## ✨ Features / Öne Çıkan Özellikler

- **[EN] Model Presets:** Instantly select top architectures like Llama-3-70B, DeepSeek-V3 (MoE), Qwen-2.5-72B, and Mistral-Large.  
  **[TR] Hazır Model Presetleri:** Llama-3-70B, DeepSeek-V3, Qwen-2.5-72B ve Mistral-Large gibi dev modelleri tek tıkla seçebilme.

- **[EN] Multi-Engine Support:** Optimize commands for **vLLM** (High-Throughput) or **SGLang** (RadixAttention & Low-Latency).  
  **[TR] Multi-Engine Desteği:** İster **vLLM** ister **SGLang** motoruyla optimize komutlar üretme.

- **[EN] Smart TP & PP Matrix:** Automatically calculate Tensor Parallelism (TP) and Pipeline Parallelism (PP) based on GPU count and context window.  
  **[TR] Akıllı TP & PP Dağılım Matrisi:** GPU sayısı ve context uzunluğuna göre TP ve PP değerlerini otomatik hesaplama.

- **[EN] Cloud Cost Simulation:** Real-time multi-GPU cluster economics and cloud cost estimation.  
  **[TR] Bulut Maliyet Simülasyonu:** Çoklu GPU cluster ekonomisini anlık hesaplama ve maliyet optimizasyonu.

- **[EN] Kubernetes (K8s) Manifest Generator:** Export production-ready K8s deployment manifests instantly.  
  **[TR] Kubernetes Manifest Üreticisi:** Üretilen konfigürasyonu doğrudan kurumsal K8s deployment manifestosuna dönüştürme.

- **[EN] Dual Language UI (i18n):** Instant switching between English and Turkish.  
  **[TR] Çok Dilli Web Paneli (i18n):** Türkçe ve İngilizce anlık dil desteği.

---

## 🛠️ Modular Architecture / Modüler Mimari
- `main.go`: Application entry point / Sunucu tetikleyici ana giriş.
- `parser.go`: VRAM, memory and cost optimization engine / Hesaplama ve optimizasyon motoru.
- `server.go`: Interactive HTTP web server, i18n & enterprise UI layer / Web sunucusu ve arayüz katmanı.

---

## 🚀 Quick Start / Çalıştırma
To run the project locally / Projeyi yerel ortamınızda ayağa kaldırmak için:
```bash
go run main.go parser.go server.go
