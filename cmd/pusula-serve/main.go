package main

import (
	"encoding/json"
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
	model := flag.String("model", "llama-3.3-70b", "model id or Hugging Face path")
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
		fmt.Printf("decode    %.0f tok/s HBM ceiling\n", a.DecodeTokPerS)
		fmt.Printf("cost      $%.2f/hr  ~$%.0f/mo  (%s %s x%d)\n", a.HourlyUSD, a.MonthlyUSD, *provider, a.GPU.ID, *gpus)
		fmt.Printf("note      %s\n\n", a.Recommendation)
		fmt.Println(a.EngineCommand)
		return
	}

	logger := slog.Default()
	mux := http.NewServeMux()
	mux.Handle("/", httpapi.New().Handler())
	mux.HandleFunc("/ops/status", HandleOpsStatus)
	mux.HandleFunc("/card", engine.CardMiddleware(logger, func(w http.ResponseWriter, r *http.Request) {
		if err := engine.GlobalQueue.AcquireSlot(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		defer engine.GlobalQueue.ReleaseSlot()

		processor := r.Context().Value("cardProcessor").(*engine.CardProcessor)
		_ = processor.AddSequenceStand(r.Context())
		_ = processor.AddPage(r.Context())
		cfg := engine.GlobalConfigManager.Get()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("pusula-serve active. model=%s max_seqs=%d (estimates only)\n", cfg.ModelName, cfg.MaxSeqs)))
	}))
	engine.RegisterTelemetryRoutes(mux)

	log.Printf("pusula-serve listening on %s — console https://pusulainfra.github.io/pusula-serve/", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func HandleOpsStatus(w http.ResponseWriter, _ *http.Request) {
	cfg := engine.GlobalConfigManager.Get()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":            "ok",
		"product":           "pusula-serve",
		"model":             cfg.ModelName,
		"max_seqs":          cfg.MaxSeqs,
		"max_model_len":     cfg.MaxModelLen,
		"stands":            engine.GlobalMetrics.ActiveStands.Load(),
		"pages":             engine.GlobalMetrics.ActivePages.Load(),
		"errors":            engine.GlobalMetrics.TotalErrors.Load(),
		"console":           "https://pusulainfra.github.io/pusula-serve/",
		"disclaimer":        "Estimates only. Not a quote, SLA, or capacity guarantee.",
	})
}
