<div align="center">
  # Pusula Serve 🚀
  *Realistic LLM Serving Config & Optimization Copilot*

  <p><b>[EN]</b> An open-source Go tool for vLLM & SGLang deployment estimation: model preset + context + GPU count -> realistic VRAM/KV cache calculation, TP/PP selection, and launch commands.</p>
  <p><b>[TR]</b> vLLM & SGLang sunumları için gerçekçi bellek tahmini, TP/PP seçimi ve başlatma komutu üreten Go tabanlı açık kaynak araç.</p>

  <p>
    <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Engines-vLLM%20%7C%20SGLang-orange?style=flat-square" alt="Engines">
    <img src="https://img.shields.io/badge/Status-MVP-blue?style=flat-square" alt="Status">
    <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  </p>
</div>

---

## 💡 What it does / Ne Yapar?
- **MVP Focus:** Model preset + context + GPU count -> VRAM/KV estimation, TP/PP, launch command.
- *Not yet included:* Traffic replay, prefix-cache hit estimation, precise GPU pricing.

---

## 🚀 Quick Start / Çalıştırma
```bash
go test ./...
go run main.go parser.go server.go

