<div align="center">

# Pusula Serve

**The chip is not the bill.**

Heterogeneous inference fleet planner for vLLM, SGLang, and mixed boxes.

[![Live](https://img.shields.io/badge/console-pusulainfra.github.io-6ee0ff)](https://pusulainfra.github.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![X](https://img.shields.io/badge/X-@pusulainfra-000)](https://x.com/pusulainfra)

</div>

**Live v1:** [https://pusulainfra.github.io/](https://pusulainfra.github.io/)

Not a VRAM toy. A serving plan.

- Fit / OOM on H100, 5090, RTX Pro, Spark GB10 UMA, Mac Studio / Mini
- GQA and MLA KV, MoE loaded vs active
- Concurrency board 1 / 4 / 16 / 32
- Prefix hit as a lever, not a slogan
- Infer vs agent role on a mixed fleet
- Prefill/decode split call (yes / no / maybe)
- Multi-LoRA tax, speculative-decode tax, ISL/OSL shape
- Launch command for vLLM / SGLang / MLX
- Scenes from production posts: prefix 0%, MoE+256K, pull latest, Spark leak

Change the flag before you buy another GPU.

## Run locally

```bash
go test ./...
go run ./cmd/pusula-serve
```

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra)
