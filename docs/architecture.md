# PusulaInfra Architecture Overview

PusulaInfra, LLM servis altyapıları için donanım gereksinimlerini, bant genişliği sınırlarını (HBM-bound) ve maliyet optimizasyonunu hesaplayan yüksek performanslı bir Go motorudur.

## Temel Bileşenler
- **Engine (`engine/`):** VRAM, KV Cache, TP/PP topolojisi ve maliyet hesaplamalarının yapıldığı çekirdek mantık katmanı.
- **HTTP API (`internal/httpapi/`):** Küme durumunu ve entegrasyon uç noktalarını sunan HTTP sunucusu.
- **CLI (`cmd/pusula-serve/`):** Komut satırı üzerinden hızlı analiz ve konfigürasyon export (JSON, YAML, Helm, Terraform) imkanı.
