<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.22-blue?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Framework-vLLM%20%2F%20SGLang-orange?style=for-the-badge&logo=ai" alt="vLLM">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker" alt="Docker">
</p>

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

## 📊 Supported Hardware & Performance Matrix

| GPU Modeli | VRAM | Bant Genişliği (TB/s) | Öngörülen Durum |
| :--- | :---: | :---: | :--- |
| **NVIDIA H100** | 80 GB | 3.35 TB/s | HBM-Bound (Mükemmel) |
| **NVIDIA A100** | 80 GB | 2.03 TB/s | Kararlı / Dengeli |
| **NVIDIA RTX 5090** | 32 GB | 1.79 TB/s | Yüksek Concurrency Sınırı |

---

## 📦 Installation & Quick Start

Clone the repository and build the binary:

```bash
git clone [https://github.com/PusulaInfra/pusula-serve.git](https://github.com/PusulaInfra/pusula-serve.git)
cd pusula-serve
go build -o pusula-serve ./cmd/pusula-serve
