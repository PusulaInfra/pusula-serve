<div align="center">

<img src="https://pusulainfra.github.io/banner.svg" alt="Pusula Serve" width="100%"/>

# Pusula Serve

**The chip is not the bill.**

LLM serving config copilot for vLLM, SGLang, KV, TP/PP.

[![Live](https://img.shields.io/badge/console-pusulainfra.github.io-6ee0ff)](https://pusulainfra.github.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![X](https://img.shields.io/badge/X-@pusulainfra-000)](https://x.com/pusulainfra)

</div>

**Live:** [https://pusulainfra.github.io/](https://pusulainfra.github.io/)

Not a VRAM toy. A serving plan.

- Fit / OOM on H100, 5090, RTX Pro, Spark GB10 UMA, Mac Studio / Mini
- GQA and MLA KV, MoE loaded vs active
- Concurrency board 1 / 4 / 16 / 32
- Prefix hit as a lever
- Infer vs agent on a mixed fleet
- P/D split call, LoRA tax, spec tax
- Launch line for vLLM / SGLang / MLX
- Spark stand: 256K + util 0.78

Change the flag before you buy another GPU.

## Run locally

```bash
go test ./...
go run ./cmd/pusula-serve
```

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra)
