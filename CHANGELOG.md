# Changelog

All notable changes to Pusula Serve are listed here.

## [1.2.0] - 2026-09-03

### Product
- User Pages (`https://pusulainfra.github.io/`) and Project Pages (`https://pusulainfra.github.io/pusula-serve/`) now ship the same interactive console: old `pack()` / `paint()` math inside the glass shell.
- Lab LED is computed from the 16/32 board (not hardcoded).
- Default preset: Llama 3.3 70B · 4×H100 · 16K · 16. `llama-3.1-70b` remains an alias.
- MINI32 VRAM is 32 GB. Ops paints on load and accepts Serve query strings.
- Disclaimer / Privacy / Terms linked from the console.

### Engine
- `HandleOpsStatus` implements `GET /ops/status`.
- CLI default model is `llama-3.3-70b`.
- `make run` aliases `make run-cli`. `run-cli` / `run-server` remain.
- `round` helper added so `engine/bandwith.go` compiles. The misspelled filename is kept.
- Decode tok/s filled from HBM bandwidth / active weights on GPU.

### Legal
- Estimates are not quotes, SLAs, or capacity guarantees.

## [1.1.0] - 2026-09-02

### Product
- GitHub Pages console is now the live planner (same engine as the Go CLI), not a static mock.
- Canonical links: GitHub `PusulaInfra/pusula-serve`, X `@pusulainfra`, console `pusulainfra.github.io/pusula-serve`.
- Disclaimer, privacy, and terms pages. LICENSE copyright is PusulaInfra.
- Hardware matrix includes L40S, RTX 5090, A100, H100, H200, and B200.

### Engine
- SLA TPOT is per-token (was incorrectly `genLen × 12.5`).
- Decode tokens/s filled from HBM bandwidth / active weights on GPU.
- Missing `round` helper, `HandleOpsStatus`, embed static, and `sync.Pool` syntax are fixed so the binary builds.

### Legal
- Estimates are not quotes, SLAs, or capacity guarantees.
- Trademark / affiliation notice. KVKK/GDPR: no product database, localStorage + URL only.

## [1.0.0] - 2026-09-02

- Initial CLI/API planner: vLLM/SGLang VRAM + KV, Green AI, Prometheus/Grafana templates, Kubernetes manifests, Docker Compose, CI.
