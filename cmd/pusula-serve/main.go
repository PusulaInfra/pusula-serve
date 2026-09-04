package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/PusulaInfra/pusula-serve/engine"
	"github.com/PusulaInfra/pusula-serve/internal/httpapi"
	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

func main() {
	fs := flag.NewFlagSet("pusula-serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { printHelp(fs) }
	addr := fs.String("addr", ":8080", "HTTP listen address")
	cli := fs.Bool("cli", false, "deprecated: same as plan")
	jsonOut := fs.Bool("json", false, "print plan as JSON")
	model := fs.String("model", "llama-3.3-70b", "model id or Hugging Face path")
	ctx := fs.Int("ctx", 16384, "max model length")
	gpus := fs.Int("gpus", 4, "GPU count")
	seqs := fs.Int("seqs", 16, "max concurrent sequences")
	gpu := fs.String("gpu", "H100", "GPU type")
	dtype := fs.Int("dtype", 2, "bytes per param: 2=bf16/fp16, 1=fp8")
	engineType := fs.String("engine", "vllm", "vllm or sglang")
	provider := fs.String("provider", "lambda", "lambda, runpod, aws, gcp")

	args := os.Args[1:]
	cmd := "serve"
	if len(args) > 0 && (args[0] == "help" || args[0] == "-h" || args[0] == "--help") {
		printHelp(fs)
		os.Exit(0)
	}
	if len(args) > 0 && !isFlag(args[0]) {
		cmd = args[0]
		args = args[1:]
	}
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Deprecated -cli always means plan, even with no subcommand.
	if *cli && cmd != "plan" && cmd != "cli" {
		cmd = "plan"
	}

	switch cmd {
	case "plan", "cli":
		os.Exit(runPlan(*model, *ctx, *gpus, *seqs, *gpu, *dtype, *engineType, *provider, *jsonOut))
	case "serve", "server":
		runServer(*addr)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printHelp(fs)
		os.Exit(1)
	}
}

func isFlag(s string) bool {
	return len(s) > 0 && s[0] == '-'
}

func runPlan(model string, ctx, gpus, seqs int, gpu string, dtype int, engineType, provider string, jsonOut bool) int {
	a := serve.Analyze(serve.ServingConfig{
		ModelName:   model,
		MaxModelLen: ctx,
		NumGpus:     gpus,
		MaxNumSeqs:  seqs,
		DtypeBytes:  dtype,
		GpuType:     gpu,
		Engine:      engineType,
		Provider:    provider,
	})

	if jsonOut {
		_ = json.NewEncoder(os.Stdout).Encode(a)
		if a.OOM {
			return 2
		}
		return 0
	}

	decision := "STANDS"
	code := 0
	if a.OOM {
		decision = "PAGES"
		code = 2
	}
	fmt.Printf("decision  %s\n", decision)
	fmt.Printf("model     %s (%s)\n", a.Spec.Name, a.Spec.HF)
	fmt.Printf("weights   %.1f GB\n", a.WeightGB)
	fmt.Printf("kv        %.1f GB\n", a.KVGB)
	fmt.Printf("per-gpu   %.1f / %.1f GB usable\n", a.PerGPUGB, a.UsablePerGPU)
	fmt.Printf("topology  TP=%d PP=%d  oom=%v\n", a.TP, a.PP, a.OOM)
	fmt.Printf("note      %s\n\n", a.Recommendation)
	fmt.Println(a.EngineCommand)
	fmt.Println("disclaimer  Estimates only. Not a quote, SLA, or capacity guarantee.")
	return code
}

func runServer(addr string) {
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

	log.Printf("pusula-serve listening on %s — console https://pusulainfra.github.io/pusula-serve/", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func HandleOpsStatus(w http.ResponseWriter, _ *http.Request) {
	cfg := engine.GlobalConfigManager.Get()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":        "ok",
		"product":       "pusula-serve",
		"model":         cfg.ModelName,
		"max_seqs":      cfg.MaxSeqs,
		"max_model_len": cfg.MaxModelLen,
		"stands":        engine.GlobalMetrics.ActiveStands.Load(),
		"pages":         engine.GlobalMetrics.ActivePages.Load(),
		"errors":        engine.GlobalMetrics.TotalErrors.Load(),
		"console":       "https://pusulainfra.github.io/pusula-serve/",
		"disclaimer":    "Estimates only. Not a quote, SLA, or capacity guarantee.",
	})
}

func printHelp(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "pusula-serve — generates a serving plan. Not a bill or SLA.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  pusula-serve plan [flags]")
	fmt.Fprintln(os.Stderr, "  pusula-serve serve [flags]")
	fmt.Fprintln(os.Stderr, "  pusula-serve            (same as serve, HTTP :8080)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Exit codes: 0 fits (STANDS), 2 OOM (PAGES), 1 user error")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  pusula-serve plan --model llama-3.3-70b --gpu H100 --gpus 4 --ctx 16384 --seqs 16")
	fmt.Fprintln(os.Stderr, "  pusula-serve serve --addr :8080")
	fmt.Fprintln(os.Stderr, "")
	fs.PrintDefaults()
}
