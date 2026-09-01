# Değişiklik Günlüğü (Changelog)

PusulaInfra projesindeki tüm önemli değişiklikler bu dosyada kronolojik olarak listelenir.

## [1.0.0] - 2026-09-02
### Eklenenler
- vLLM ve SGLang motorları için VRAM ve KV Cache optimizasyon modülleri eklendi.
- Green AI (TDP güç tüketimi ve karbon salınımı hesaplama) özellikleri entegre edildi.
- Prometheus metrikleri ve enterprise alert yapılandırmaları (`ops/prometheus-alerts.yaml`) eklendi.
- Grafana monitoring dashboard şablonu (`ops/grafana-dashboard.json`) eklendi.
- Kubernetes deployment manifestoları ve Docker Compose altyapısı hazırlandı.
- Otomatik CI/CD (`ci.yml`) ve Makefile otomasyon araçları kuruldu.
