<div align="center">

<img src="https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=1400&q=80" alt="Pusula Serve" width="100%">

# Pusula Serve

**[EN]** LLM serving config copilot for vLLM & SGLang.  
**[TR]** vLLM ve SGLang için serving ayarı, bellek ve KV tahmini.

Model + context + GPU → VRAM · KV cache · TP/PP · launch command

<br/>

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![vLLM](https://img.shields.io/badge/vLLM-supported-orange?style=for-the-badge)](https://github.com/vllm-project/vllm)
[![SGLang](https://img.shields.io/badge/SGLang-supported-orange?style=for-the-badge)](https://github.com/sgl-project/sglang)
[![Status](https://img.shields.io/badge/Status-MVP-blue?style=for-the-badge)](#)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

[English](#english) · [Türkçe](#türkçe) · [API](#api) · [Roadmap](#roadmap)

</div>

---

<a id="english"></a>

## English

Open-source Go tool. You pick a model, context length and GPU count. It estimates **weight memory**, **KV cache**, **per-GPU use**, **OOM risk**, then prints a **vLLM / SGLang** command and a rough Kubernetes snippet.

It does **not** replay production traffic. It does **not** quote a real cloud invoice. Config first.

### Why it exists

Teams often buy another GPU when `max-model-len` or `max-num-seqs` is the bill. KV grows with context × concurrency. This tool makes that visible before you launch.

### Features

| Feature | Status |
|---|---|
| Model presets (Llama 70B, Qwen 72B, Mistral Large, DeepSeek-V3 MoE) | Done |
| Weight vs KV split | Done |
| TP / PP suggestion | Done |
| OOM flag (80GB GPU, 0.90 util default) | Done |
| vLLM + SGLang launch command | Done |
| JSON API | Done |
| Traffic replay / prefix-cache hit | Not yet |
| Exact cloud pricing | Not yet |

### Quick start

```bash
go test ./...
go run main.go parser.go server.go