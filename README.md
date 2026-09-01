<div align="center">

  <h1>Pusula Serve</h1>
  <p><strong>Enterprise-Grade LLM Serving Engine & Cluster Capacity Optimizer</strong></p>
  
  <p>
    <a href="#-architecture">Architecture</a> •
    <a href="#-features">Features</a> •
    <a href="#-quickstart">Quickstart</a> •
    <a href="#-benchmarks">Benchmarks</a>
  </p>

  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Engine-vLLM%20%7C%20SGLang-orange?style=for-the-badge" alt="Engine Support">
  <img src="https://img.shields.io/badge/CI%2FRace-Passing-success?style=for-the-badge&logo=githubactions" alt="CI Status">

  <p><em>"Not a VRAM toy. A serving plan. The chip is not the bill."</em></p>
</div>

---

## ⚡ Overview

**Pusula Serve** is a high-performance, production-ready inference orchestration and capacity planning engine written in Go. It bridges the gap between hardware capability and real-world LLM serving constraints, featuring strict sequence slot enforcement, zero-allocation memory pooling, and profile-guided optimization (PGO).

---

## 🏗️ Core Architecture

* **Engine Core (`/engine`)**: Modular concurrent management handling precise KV cache allocation, state transitions, and high-precision telemetry.
* **Concurrency Guard**: Strict enforcement of **16 Sequence Stands** and **32 Sequence Pages** with asynchronous rate-limiting and graceful queue fallback.
* **Zero-Allocation Pools**: Leverages `sync.Pool` to eliminate garbage collection pressure during peak token generation spikes.
* **PGO Optimized**: Built with Go Profile-Guided Optimization (`-pgo=auto`) for maximum CPU throughput.

---

## 🚀 Quickstart

Clone the repository and spin up the enterprise engine using the automated Makefile:

```bash
# Clone the repository
git clone [https://github.com/PusulaInfra/pusula-serve.git](https://github.com/PusulaInfra/pusula-serve.git)
cd pusula-serve

# Run tests with race detector
make test

# Build with Profile-Guided Optimization
make build

# Start the enterprise server
make run
