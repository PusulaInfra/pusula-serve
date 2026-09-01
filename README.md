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

## Estimate vs Measure

Estimate is the plan. It does not touch a GPU.

Measure is evidence from this box, this run.

`Measured on this box, this run. Not a vendor SLA.`

A measured tok/s is never a PAGES guarantee. If `nvidia-smi` is missing, Measure says so. It does not invent throughput.

```bash
go test ./...
go run ./cmd/pusula-serve --cli --model llama-3.3-70b --gpu H100 --gpus 4 --ctx 8192 --seqs 8
go run ./cmd/pusula-serve --measure --bench-sec 5
go run ./cmd/pusula-serve --live-vram
go run ./cmd/pusula-serve --apply=dry-run
# exec only with both flags:
go run ./cmd/pusula-serve --apply=exec --remote
```

HTTP:

- `POST /api/analyze` — estimate only (contract unchanged)
- `POST /api/measure` — estimate + live VRAM + bench skip
- `GET /api/live-vram`
- `POST /api/apply` — dry-run unless `apply=exec` and `remote=true`

Change the flag before you buy another GPU.

## Brand

[x.com/pusulainfra](https://x.com/pusulainfra)
