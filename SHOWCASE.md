# Pusula Serve: Enterprise-Grade LLM Serving Engine in Go

We built **Pusula Serve** because most LLM VRAM calculators are toys that ignore real-world serving constraints. When deploying 70B+ models on H100 clusters, hardware limits, memory pooling, and concurrency guards dictate whether your system survives production or hits OOM errors.

## Why Go?
We chose Go to bypass Python/Ray runtime overhead. By implementing:
- **Strict Concurrency Guards**: Hard enforcement of 16 Sequence Stands and 32 Sequence Pages.
- **Zero-Allocation Memory Pooling**: Using `sync.Pool` to eliminate GC pressure during heavy token generation.
- **Profile-Guided Optimization (PGO)**: Compiled with `-pgo=auto` for maximum CPU throughput.
- **Native Observability**: Prometheus-ready metrics and Kubernetes deployment readiness out of the box.

Check out the interactive console and production manifests on [GitHub](https://github.com/PusulaInfra/pusula-serve).
