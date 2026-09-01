# PusulaInfra 🧭 Mimari ve Akış Rehberi

Bu doküman, `PusulaInfra` motorunun donanım analizi, bellek yönetimi ve operasyonel akışını görsel ve şematik olarak açıklamaktadır.

---

## 🔄 1. İstek ve Analiz Akış Şeması (CLI & API)

Kullanıcı bir analiz başlattığında veya sunucuya istek attığında verinin izlediği yol şu şekildedir:

```text
[ Kullanıcı / CLI / Browser ]
          │
          ├──> 1. CLI Modu (`-cli`) ──> [ serve.Analyze() ] ──> Terminal Çıktısı (Tablo/Rapor)
          │
          └──> 2. HTTP Server (`:8080`) 
                    │
                    ├──> /ops/status ──> [ OpsStatusResponse ] (Sağlık & VRAM Durumu)
                    ├──> /card       ──> [ GlobalQueue / Rate Limiter ] ──> Koruma & Slot Kontrolü
                    └──> /metrics    ──> [ Telemetry Routes ] (Prometheus Metrikleri)
