<div align="center">

# Pusula Serve

**The chip is not the bill.**

Production planner for vLLM and SGLang.

[![Live](https://img.shields.io/badge/console-pusulainfra.github.io-6ee0ff)](https://pusulainfra.github.io/)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![X](https://img.shields.io/badge/X-@pusulainfra-000)](https://x.com/pusulainfra)

</div>

**Open now:** [https://pusulainfra.github.io/](https://pusulainfra.github.io/)

## What operators actually need

Generic calculators stop at `params × bytes`.

| Question | Generic VRAM toy | Pusula Serve |
|---|---|---|
| Weights | params × dtype | loaded weights |
| KV | often GQA-only | GQA **and MLA** |
| MoE | active params look cheap | loaded experts stay on HBM |
| Concurrency | slider | max-num-seqs that still fits |
| Prefix cache | missing | hit-rate lever |
| Output | a GB number | vLLM / SGLang command + K8s |
| Honesty | "required VRAM" | estimate, not a profiler |

Change the flag before you buy another GPU.

## Run locally

```bash
go test ./...
go run ./cmd/pusula-serve
# http://localhost:8080
```

```bash
go run ./cmd/pusula-serve -cli -model deepseek-v3 -ctx 65536 -gpus 8 -gpu H100 -engine vllm
```

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra) — what actually works in prod.
