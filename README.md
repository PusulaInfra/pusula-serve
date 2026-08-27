<div align="center">

# Pusula Serve 🚀
*Realistic LLM Serving Config & Optimization Copilot*

[![Language](https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![Engines](https://img.shields.io/badge/Engines-vLLM%20%7C%20SGLang-orange?style=flat-square)](https://github.com/vllm-project/vllm)
[![Status](https://img.shields.io/badge/Status-MVP-blue?style=flat-square)]()
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

</div>

---

> **[EN]** An open-source Go tool for vLLM & SGLang deployment estimation: model preset + context + GPU count -> realistic VRAM/KV cache calculation, TP/PP selection, and launch commands.
> 
> **[TR]** vLLM & SGLang sunumları için gerçekçi bellek tahmini, TP/PP seçimi ve başlatma komutu üreten Go tabanlı açık kaynak araç.

---

## 🖼️ Preview / Arayüz Önizlemesi

> *Aşağıdaki görselin görünmesi için projenin içine ekran görüntüsü ekleyip yolunu güncelleyebilirsin (Örn: `docs/preview.png`).*

<div align="center">
  <img src="https://via.placeholder.com/900x450/1e1e1e/00ADD8?text=Pusula+Serve+UI+Preview" alt="Pusula Serve UI Preview" width="100%">
</div>

---

## 💡 What it does / Ne Yapar?

* **MVP Focus:** Model preset + context + GPU count + concurrency -> VRAM/KV cache estimation, TP/PP selection, launch command.
* **Not yet included:** Traffic replay, prefix-cache hit estimation, precise GPU pricing.

---

## 🚀 Quick Start / Çalıştırma

Gereksinimleri kurduktan sonra projeyi hızlıca test etmek ve çalıştırmak için şu komutları kullanabilirsin:

```bash
# Testleri çalıştır
go test ./...

# Projeyi başlat
go run main.go parser.go server.go
