package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PusulaInfra/pusula-serve/internal/apply"
	"github.com/PusulaInfra/pusula-serve/internal/httpapi"
	"github.com/PusulaInfra/pusula-serve/internal/live"
	"github.com/PusulaInfra/pusula-serve/internal/measure"
	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	cli := flag.Bool("cli", false, "print analysis to stdout and exit")
	measureFlag := flag.Bool("measure", false, "print measure report (estimate + live VRAM + bench skip)")
	liveVRAM := flag.Bool("live-vram", false, "print live VRAM snapshot and exit")
	benchSec := flag.Int("bench-sec", 5, "bench window seconds (no invented tok/s)")
	applyMode := flag.String("apply", "", "dry-run or exec")
	remote := flag.Bool("remote", false, "required together with --apply=exec")
	model := flag.String("model", "llama-3.3-70b", "model id or Hugging Face path")
	ctx := flag.Int("ctx", 16384, "max model length")
	gpus := flag.Int("gpus", 4, "GPU count")
	seqs := flag.Int("seqs", 16, "max concurrent sequences")
	gpu := flag.String("gpu", "H100", "GPU type")
	dtype := flag.Int("dtype", 2, "bytes per param: 2=bf16/fp16, 1=fp8")
	engine := flag.String("engine", "vllm", "vllm or sglang")
	provider := flag.String("provider", "lambda", "lambda, runpod, aws, gcp")
	flag.Parse()

	cfg := serve.ServingConfig{
		ModelName:   *model,
		MaxModelLen: *ctx,
		NumGpus:     *gpus,
		MaxNumSeqs:  *seqs,
		DtypeBytes:  *dtype,
		GpuType:     *gpu,
		Engine:      *engine,
		Provider:    *provider,
	}

	if *liveVRAM {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(live.SnapshotNow())
		return
	}

	if *measureFlag {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(measure.Run(cfg, *benchSec))
		return
	}

	if *applyMode != "" {
		a := serve.Analyze(cfg)
		res, err := apply.Run(apply.Request{Mode: *applyMode, Remote: *remote, Line: a.EngineCommand})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("apply     %s ran=%v\n", res.Mode, res.Ran)
		fmt.Printf("note      %s\n", res.Note)
		fmt.Println(res.Line)
		return
	}

	if *cli {
		a := serve.Analyze(cfg)
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

	log.Printf("pusula-serve  http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, httpapi.New().Handler()); err != nil {
		log.Fatal(err)
	}
}
