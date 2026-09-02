<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.22-blue?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Console-GitHub%20Pages-222?style=for-the-badge" alt="Pages">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/X-@pusulainfra-000?style=for-the-badge" alt="X">
</p>

# Pusula Serve

> **Not a VRAM toy. A serving plan. The chip is not the bill.**

One product: a Go CLI/API and a live web console that share the same serving math. Plan vLLM / SGLang deployments — weights, KV, TP/PP, cloud list-price, HBM decode ceiling, SLA, LoRA, MoE/MLA — before you buy the box.

**Live console:** https://pusulainfra.github.io/pusula-serve/  
**Source:** https://github.com/PusulaInfra/pusula-serve  
**X:** https://x.com/pusulainfra

Default walk-through: **Llama 3.3 70B · 4× H100 · 16K · 16 seqs**.

**STANDS** fits safely. **PAGES** trigger OOMs. Guard: 16 sequence stands, 32 sequence pages.

> Estimates only — not a quote, SLA, or capacity guarantee. See [DISCLAIMER.md](DISCLAIMER.md).

---

## Features

- GQA and MLA KV (DeepSeek is not sized as GQA)
- Quant matrix: FP16, FP8, AWQ, GPTQ, GGUF
- Multi-LoRA VRAM reserve
- HBM-bound decode ceiling
- TTFT / TPOT SLA (TPOT is per-token)
- Green AI: TDP, kWh, CO₂, electricity
- Resilience, P/D split, interconnect notes
- Export: vLLM, SGLang, Kubernetes, Helm, Terraform, JSON
- Shareable plan URL (console)

---

## Hardware matrix

| GPU | VRAM | Bandwidth | Notes |
| :--- | :---: | :---: | :--- |
| NVIDIA L40S | 48 GB | 0.84 TB/s | Inference / adapter boxes |
| NVIDIA RTX 5090 | 32 GB | 1.79 TB/s | High concurrency limit |
| NVIDIA A100 | 80 GB | 2.03 TB/s | Balanced |
| NVIDIA H100 | 80 GB | 3.35 TB/s | Default showcase |
| NVIDIA H200 | 141 GB | 4.80 TB/s | Long context |
| NVIDIA B200 | 192 GB | 8.00 TB/s | 405B-class |

Cloud rate cards: Lambda, RunPod, AWS, GCP — **list-price snapshots**, not invoices.

---

## Quick start

```bash
git clone https://github.com/PusulaInfra/pusula-serve.git
cd pusula-serve
go build -o pusula-serve ./cmd/pusula-serve
./pusula-serve -cli -model llama-3.3-70b -gpu H100 -gpus 4 -ctx 16384 -seqs 16 -provider lambda
```

HTTP API + landing (same engine):

```bash
./pusula-serve -addr :8080
# POST /api/analyze
# GET  /api/models
# GET  /api/health
```

Make targets: `make build`, `make test`, `make run-cli`.

---

## Console

The official installable web console is GitHub Pages:

**https://pusulainfra.github.io/pusula-serve/**

Add it to the home screen from that origin (standalone PWA). The `docs/` folder **is** the console — not a screenshot of it.

---

## Legal

- [Disclaimer](DISCLAIMER.md) · [Privacy](PRIVACY.md) · [Terms](TERMS.md) · [Security](SECURITY.md) · [NOTICE](NOTICE)
- MIT License — copyright PusulaInfra
- Not affiliated with NVIDIA, Meta, Mistral, DeepSeek, vLLM, SGLang, AWS, GCP, Lambda, or RunPod
- Launch scripts may name Hugging Face repos; you must accept each model license

---

## Links

| | |
| --- | --- |
| Console | https://pusulainfra.github.io/pusula-serve/ |
| GitHub | https://github.com/PusulaInfra/pusula-serve |
| X | https://x.com/pusulainfra |
