<div align="center">

# Pusula Serve

**The chip is not the bill.**

Production planner for vLLM and SGLang.

Open the console: [`docs/index.html`](docs/index.html)

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![X](https://img.shields.io/badge/X-@pusulainfra-000)](https://x.com/pusulainfra)

</div>

## What operators actually need

Generic calculators stop at `params × bytes`.

Serving bills break after that:

| Question | Generic VRAM toy | Pusula Serve |
|---|---|---|
| Weights | params × dtype | same, labeled as **loaded** weights |
| KV | often GQA-only | **GQA and MLA** (DeepSeek latent+RoPE) |
| MoE | active params look cheap | loaded experts stay on HBM |
| Concurrency | ignored or a slider | **max-num-seqs that still fits** |
| Prefix cache | missing | hit-rate lever on KV |
| Output | a GB number | **vLLM / SGLang command + K8s sketch** |
| Honesty | "required VRAM" | estimate, not a profiler |

Change the flag before you buy another GPU.

## Console

```bash
go test ./...
go run ./cmd/pusula-serve
# http://localhost:8080
```

Or open `docs/index.html` — no server required.

GitHub Pages: Settings → Pages → Deploy from GitHub Actions (workflow already in-repo).
Live URL after first successful Actions run:
`https://pusulainfra.github.io/pusula-serve/`

## CLI

```bash
go run ./cmd/pusula-serve -cli \
  -model deepseek-v3 \
  -ctx 65536 \
  -gpus 8 \
  -gpu H100 \
  -engine vllm
```

## API

`GET /api/health`  
`GET /api/models`  
`POST /api/analyze`

## Not this product

It does not replay production traces.
Cloud dollars are list-price sketches.
It will not tell you p99 under real traffic.

It will tell you when **context × concurrency** is the invoice, not the chip.

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra) — what actually works in prod.
