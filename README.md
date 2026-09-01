# PusulaInfra 🧭

> **Enterprise-Grade LLM Infrastructure, Cost, and Operational Analysis Engine**

PusulaInfra is a high-performance, Go-based backend and CLI tool designed for AI infrastructure engineers, ML platform teams, and cloud architects. It calculates precise GPU memory requirements, KV cache scaling, network bandwidth bounds (HBM-bound), and multi-cloud operational costs for modern LLM serving frameworks (vLLM & SGLang).

---

## 🚀 Key Features & Capabilities

- **🧠 Deep Memory & VRAM Modeling**: Precise calculation for model weights, KV cache scaling, and activation overheads across different quantization formats (FP16, FP8, AWQ, GPTQ, GGUF).
- **⚡ Bandwidth & Decode Tapping**: Calculates maximum theoretical token throughput based on GPU memory bandwidth (HBM-bound limits).
- **🎛️ Multi-LoRA Memory Overhead**: Dynamic VRAM estimation for simultaneous fine-tuned adapter scaling.
- **🌱 Green AI & Carbon Footprint**: Real-time TDP power consumption estimation, monthly electricity costs, and carbon emission analytics.
- **⏱️ Latency SLA Simulation**: Predicts Time-to-First-Token (TTFT) and Token-Per-Output-Token (TPOT) under various concurrency loads.
- **🛡️ Resilience & Failover Modeling**: Multi-node cluster redundancy and capacity drop analysis during node outages.
- **📦 Multi-Format Infrastructure Export**: Instantly export serving configurations to **JSON, YAML, Kubernetes ConfigMaps, Helm Values, and Terraform snippets**.
- **🚦 Enterprise Rate Limiting & Queue**: Advanced slot management and telemetry endpoints (`/ops/status`, `/metrics`) for production monitoring.

---

## 📦 Installation & Quick Start

Clone the repository and build the binary:

```bash
git clone [https://github.com/PusulaInfra/pusula-serve.git](https://github.com/PusulaInfra/pusula-serve.git)
cd pusula-serve
go build -o pusula-serve ./cmd/pusula-serve
