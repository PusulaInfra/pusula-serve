package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/PusulaInfra/pusula-serve/engine"
	"github.com/PusulaInfra/pusula-serve/internal/httpapi"
	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	cli := flag.Bool("cli", false, "print analysis to stdout and exit")
	model := flag.String("model", "llama-3.1-70b", "model id or Hugging Face path")
	ctx := flag.Int("ctx", 16384, "max model length")
	gpus := flag.Int("gpus", 4, "GPU count")
	seqs := flag.Int("seqs", 16, "max concurrent sequences")
	gpu := flag.String("gpu", "H100", "GPU type")
	dtype := flag.Int("dtype", 2, "bytes per param: 2=bf16/fp16, 1=fp8")
	engineType := flag.String("engine", "vllm", "vllm or sglang")
	provider := flag.String("provider", "lambda", "lambda, runpod, aws, gcp")
	flag.Parse()

	if *cli {
		a := serve.Analyze(serve.ServingConfig{
			ModelName:   *model,
			MaxModelLen: *ctx,
			NumGpus:     *gpus,
			MaxNumSeqs:  *seqs,
			DtypeBytes:  *dtype,
			GpuType:     *gpu,
			Engine:      *engineType,
			Provider:    *provider,
		})
		fmt.Printf("model     %s (%s)\n", a.Spec.Name, a.Spec.HF)
		fmt.Printf("weights   %.1f GB\n", a.WeightGB)
		fmt.Printf("kv        %.1f GB\n", a.KVGB)
		fmt.Printf("per-gpu   %.1f / %.1f GB usable\n", a.PerGPUGB, a.UsablePerGPU)
		fmt.Printf("topology  TP=%d PP=%d  oom=%v\n", a.TP, a.PP, a.OOM)
		fmt.Printf("cost      $%.2f/hr  ~$%.0f/mo  (%s %s x%d)\n", a.HourlyUSD, a.MonthlyUSD, *provider, a.GPU.ID, *gpus)
		fmt.Printf("note      %s\n\n", a.Recommendation)
		fmt.Println(a.EngineCommand)
		return
	}

	logger := slog.Default()
	mux := http.NewServeMux()

	// Mevcut API rotasını dahil et
	mux.Handle("/", httpapi.New().Handler())
	// Ops durum rotasını ekle
	mux.HandleFunc("/ops/status", HandleOpsStatus)

	// Gelişmiş Kuyruk (Rate Limiting) ve 16/32 Sınırlarına Sahip /card Rotası
	mux.HandleFunc("/card", engine.CardMiddleware(logger, func(w http.ResponseWriter, r *http.Request) {
		// İstek geldiğinde GlobalQueue üzerinden slot kapmaya çalışıyoruz (Aşırı yük koruması)
		if err := engine.GlobalQueue.AcquireSlot(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		defer engine.GlobalQueue.ReleaseSlot()

		processor := r.Context().Value("cardProcessor").(*engine.CardProcessor)
		requestCtx := r.Context()

		_ = processor.AddSequenceStand(requestCtx)
		_ = processor.AddPage(requestCtx)

		// Güncel Hot-Reload konfigürasyonunu okuyoruz
		cfg := engine.GlobalConfigManager.Get()

		w.Write([]byte(fmt.Sprintf("pusula-serve active! Model: %s, Max Seqs: %d", cfg.ModelName, cfg.MaxSeqs)))
	}))

	// Telemetri ve metrik rotalarını ekle (/metrics, /healthz, vb.)
	engine.RegisterTelemetryRoutes(mux)

	log.Printf("pusula-serve enterprise engine running on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
