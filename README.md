<div align="center">

# Pusula Serve

**The chip is not the bill.**

Production planner for vLLM and SGLang.

Weights · KV architecture (GQA / MLA) · TP/PP · OOM · max-num-seqs that actually fits · launch command

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![X](https://img.shields.io/badge/X-@pusulainfra-000)](https://x.com/pusulainfra)

</div>

## Why this, not another VRAM toy

Most calculators multiply `params × bytes` and stop.

Serving bills break on the next layer:

- KV grows with **context × concurrency**, not with the model card
- MoE looks cheap on active params and still loads expert weights
- DeepSeek-style **MLA** is not GQA
- Pulling `latest` changes batch and cache defaults on the same weights

Pusula Serve is a serving config console. Change the flag before you buy another GPU.

It does **not** replay production traffic. Cloud numbers are list-price sketches.

## Live UI

Open `docs/index.html` or:

```bash
go test ./...
go run ./cmd/pusula-serve
# http://localhost:8080
```

GitHub Pages: Settings → Pages → `/docs`.

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

## One repo

This repository is the product.

Archive or delete `PusulaInfra/pusula-infra`. It was a broken Wails experiment with a checked-in `.exe`.

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra) — what actually works in prod.
