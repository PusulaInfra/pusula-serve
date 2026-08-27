<div align="center">

# Pusula Serve

vLLM / SGLang serving config tahmini.

Model + context + GPU → VRAM, KV cache, TP/PP, başlatma komutu.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](#)
[![vLLM](https://img.shields.io/badge/vLLM-orange?style=flat-square)](#)
[![SGLang](https://img.shields.io/badge/SGLang-orange?style=flat-square)](#)
[![MVP](https://img.shields.io/badge/Status-MVP-blue?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

</div>

## Ne işe yarar

Serving ayarını tahmin etmeden önce bak:

- model ağırlığı (VRAM)
- KV cache (context × concurrency)
- GPU başına kullanım ve OOM riski
- TP / PP
- vLLM veya SGLang komutu
- kaba K8s Deployment taslağı

Henüz yok: gerçek traffic replay, prefix-cache hit oranı, birebir cloud fiyatı.

## Çalıştır

```bash
go test ./...
go run main.go parser.go server.go