package main

import (
	"fmt"
	"math"
	"strings"
)

type ServingConfig struct {
	ModelName            string
	MaxModelLen          int
	GpuMemoryUtilization float64
	NumGpus              int
	MaxNumSeqs           int
	DtypeBytes           int // 2=fp16/bf16, 1=fp8
}

type ModelSpec struct {
	ID           string
	ParamsB      float64 // loaded weights (total experts for MoE)
	ActiveB      float64
	Layers       int
	Hidden       int
	KVHeads      int
	HeadDim      int
	MoE          bool
}

type Analysis struct {
	Spec           ModelSpec
	WeightGB       float64
	KVGB           float64
	ActivationGB   float64
	TotalGB        float64
	PerGPUGB       float64
	UsablePerGPU   float64
	TP             int
	PP             int
	OOM            bool
	HeadroomGB     float64
	Recommendation string
}

var modelRegistry = []ModelSpec{
	{ID: "llama-3-70b", ParamsB: 70.6, ActiveB: 70.6, Layers: 80, Hidden: 8192, KVHeads: 8, HeadDim: 128, MoE: false},
	{ID: "qwen-2.5-72b", ParamsB: 72.7, ActiveB: 72.7, Layers: 80, Hidden: 8192, KVHeads: 8, HeadDim: 128, MoE: false},
	{ID: "mistral-large", ParamsB: 123, ActiveB: 123, Layers: 88, Hidden: 12288, KVHeads: 8, HeadDim: 128, MoE: false},
	{ID: "deepseek-v3", ParamsB: 671, ActiveB: 37, Layers: 61, Hidden: 7168, KVHeads: 8, HeadDim: 128, MoE: true},
	{ID: "llama-3-8b", ParamsB: 8.0, ActiveB: 8.0, Layers: 32, Hidden: 4096, KVHeads: 8, HeadDim: 128, MoE: false},
}

func ResolveModel(name string) ModelSpec {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "deepseek"):
		return modelRegistry[3]
	case strings.Contains(n, "qwen") && (strings.Contains(n, "72") || strings.Contains(n, "2.5-72")):
		return modelRegistry[1]
	case strings.Contains(n, "mistral") && strings.Contains(n, "large"):
		return modelRegistry[2]
	case strings.Contains(n, "8b"):
		return modelRegistry[4]
	case strings.Contains(n, "70"):
		return modelRegistry[0]
	default:
		s := modelRegistry[0]
		s.ID = "unknown-approx-70b"
		return s
	}
}

func pickParallelism(numGPUs, layers int, weightsGB, kvGB, gpuGB, util float64) (tp, pp int) {
	if numGPUs < 1 {
		numGPUs = 1
	}
	usable := gpuGB * util

	for tp = 1; tp <= numGPUs; tp *= 2 {
		if numGPUs%tp != 0 {
			continue
		}
		pp = numGPUs / tp
		if pp > layers {
			continue
		}
		perGPU := (weightsGB / float64(tp*pp)) + (kvGB / float64(tp))
		if perGPU <= usable {
			return tp, pp
		}
	}
	return numGPUs, 1
}

func Analyze(cfg ServingConfig) Analysis {
	if cfg.MaxModelLen < 1 {
		cfg.MaxModelLen = 8192
	}
	if cfg.NumGpus < 1 {
		cfg.NumGpus = 1
	}
	if cfg.GpuMemoryUtilization <= 0 || cfg.GpuMemoryUtilization > 0.98 {
		cfg.GpuMemoryUtilization = 0.90
	}
	if cfg.MaxNumSeqs < 1 {
		cfg.MaxNumSeqs = 16
	}
	if cfg.DtypeBytes < 1 {
		cfg.DtypeBytes = 2
	}

	spec := ResolveModel(cfg.ModelName)
	weightGB := spec.ParamsB * float64(cfg.DtypeBytes)
	
	kvBytes := 2.0 * float64(spec.Layers*spec.KVHeads*spec.HeadDim) *
		float64(cfg.MaxModelLen) * float64(cfg.MaxNumSeqs) * float64(cfg.DtypeBytes)
	kvGB := kvBytes / (1024 * 1024 * 1024)
	actGB := 2.0

	const gpuGB = 80.0
	tp, pp := pickParallelism(cfg.NumGpus, spec.Layers, weightGB, kvGB, gpuGB, cfg.GpuMemoryUtilization)

	total := weightGB + kvGB + actGB
	perGPU := (weightGB / float64(tp*pp)) + (kvGB / float64(tp)) + actGB
	usable := gpuGB * cfg.GpuMemoryUtilization
	oom := perGPU > usable

	a := Analysis{
		Spec:         spec,
		WeightGB:     weightGB,
		KVGB:         kvGB,
		ActivationGB: actGB,
		TotalGB:      total,
		PerGPUGB:     perGPU,
		UsablePerGPU: usable,
		TP:           tp,
		PP:           pp,
		OOM:          oom,
		HeadroomGB:   usable - perGPU,
	}

	switch {
	case oom:
		a.Recommendation = fmt.Sprintf(
			"OOM riski: %.1f GB/GPU > %.1f GB usable. Context'i düşür, max-num-seqs azalt, GPU ekle veya FP8 dene. MoE total weights=%.0fB, active=%.0fB.",
			perGPU, usable, spec.ParamsB, spec.ActiveB)
	case a.KVGB > a.WeightGB*0.5:
		a.Recommendation = fmt.Sprintf(
			"KV (%.1f GB) ağırlıkları yakalıyor. Fatura context/concurrency; model değil. max-model-len veya prefix-cache şart.", a.KVGB)
	case spec.MoE:
		a.Recommendation = fmt.Sprintf(
			"MoE: yüklü ağırlık %.0fB, aktif %.0fB. Token ucuz görünür; %dK context KV'yi yine şişirir.",
			spec.ParamsB, spec.ActiveB, cfg.MaxModelLen/1024)
	default:
		a.Recommendation = fmt.Sprintf("TP=%d PP=%d, headroom %.1f GB/GPU. gpu-memory-utilization=%.2f.", tp, pp, a.HeadroomGB, cfg.GpuMemoryUtilization)
	}
	return a
}

func CalculateVRAMAndCost(cfg ServingConfig) (float64, string) {
	a := Analyze(cfg)
	return a.TotalGB, a.Recommendation
}

func AnalyzeConfig(cfg ServingConfig) {
	a := Analyze(cfg)
	fmt.Printf("Model: %s (%s)\n", cfg.ModelName, a.Spec.ID)
	fmt.Printf("Weights: %.1f GB  KV: %.1f GB  Total: %.1f GB  PerGPU: %.1f GB\n", a.WeightGB, a.KVGB, a.TotalGB, a.PerGPUGB)
	fmt.Printf("TP=%d PP=%d OOM=%v\n%s\n", a.TP, a.PP, a.OOM, a.Recommendation)
}

func ClampInt(v, lo, hi int) int {
	return int(math.Min(float64(hi), math.Max(float64(lo), float64(v))))
}
