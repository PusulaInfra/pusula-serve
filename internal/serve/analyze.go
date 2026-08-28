package serve

import (
	"fmt"
	"math"
	"strings"
)

type AttnKind string

const (
	AttnGQA AttnKind = "gqa"
	AttnMLA AttnKind = "mla"
)

type ServingConfig struct {
	ModelName            string  `json:"model"`
	MaxModelLen          int     `json:"max_model_len"`
	GpuMemoryUtilization float64 `json:"gpu_memory_utilization"`
	NumGpus              int     `json:"num_gpus"`
	MaxNumSeqs           int     `json:"max_num_seqs"`
	DtypeBytes           int     `json:"dtype_bytes"`
	KVDtypeBytes         int     `json:"kv_dtype_bytes"`
	GpuType              string  `json:"gpu_type"`
	GpuVramGB            float64 `json:"gpu_vram_gb"`
	Engine               string  `json:"engine"`
	Provider             string  `json:"provider"`
	PrefixHit            float64 `json:"prefix_hit"`
}

type ModelSpec struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	HF         string   `json:"hf"`
	Family     string   `json:"family"`
	ParamsB    float64  `json:"params_b"`
	ActiveB    float64  `json:"active_b"`
	Layers     int      `json:"layers"`
	Hidden     int      `json:"hidden"`
	KVHeads    int      `json:"kv_heads"`
	HeadDim    int      `json:"head_dim"`
	MoE        bool     `json:"moe"`
	Attn       AttnKind `json:"attn"`
	MLALatent  int      `json:"mla_latent"`
	MLARopeDim int      `json:"mla_rope_dim"`
}

type GPUSpec struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	VramGB   float64 `json:"vram_gb"`
	BandGBs  float64 `json:"bandwidth_gbs"`
	TFLOPS16 float64 `json:"tflops_16"`
}

type Analysis struct {
	Spec            ModelSpec `json:"spec"`
	GPU             GPUSpec   `json:"gpu"`
	WeightGB        float64   `json:"weight_gb"`
	KVGB            float64   `json:"kv_gb"`
	KVPerTokenKB    float64   `json:"kv_per_token_kb"`
	RuntimeGB       float64   `json:"runtime_gb"`
	TotalGB         float64   `json:"total_gb"`
	PerGPUGB        float64   `json:"per_gpu_gb"`
	UsablePerGPU    float64   `json:"usable_per_gpu_gb"`
	TP              int       `json:"tp"`
	PP              int       `json:"pp"`
	OOM             bool      `json:"oom"`
	HeadroomGB      float64   `json:"headroom_gb"`
	KVShare         float64   `json:"kv_share"`
	HourlyUSD       float64   `json:"hourly_usd"`
	MonthlyUSD      float64   `json:"monthly_usd"`
	MaxFitSeqs      int       `json:"max_fit_seqs"`
	DecodeTokPerS   float64   `json:"decode_tok_per_s"`
	Recommendation  string    `json:"recommendation"`
	Fix             string    `json:"fix"`
	VLLMCommand     string    `json:"vllm_command"`
	SGLangCommand   string    `json:"sglang_command"`
	KubeSnippet     string    `json:"kube_snippet"`
	EngineCommand   string    `json:"engine_command"`
	PrefixEffective float64   `json:"prefix_effective"`
}

var ModelRegistry = []ModelSpec{
	{ID: "llama-3.1-8b", Name: "Llama 3.1 8B", HF: "meta-llama/Llama-3.1-8B-Instruct", Family: "Llama", ParamsB: 8.0, ActiveB: 8.0, Layers: 32, Hidden: 4096, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "llama-3.3-70b", Name: "Llama 3.3 70B", HF: "meta-llama/Llama-3.3-70B-Instruct", Family: "Llama", ParamsB: 70.6, ActiveB: 70.6, Layers: 80, Hidden: 8192, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "llama-3.1-405b", Name: "Llama 3.1 405B", HF: "meta-llama/Llama-3.1-405B-Instruct", Family: "Llama", ParamsB: 405, ActiveB: 405, Layers: 126, Hidden: 16384, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "qwen-2.5-7b", Name: "Qwen2.5 7B", HF: "Qwen/Qwen2.5-7B-Instruct", Family: "Qwen", ParamsB: 7.6, ActiveB: 7.6, Layers: 28, Hidden: 3584, KVHeads: 4, HeadDim: 128, Attn: AttnGQA},
	{ID: "qwen-2.5-32b", Name: "Qwen2.5 32B", HF: "Qwen/Qwen2.5-32B-Instruct", Family: "Qwen", ParamsB: 32.5, ActiveB: 32.5, Layers: 64, Hidden: 5120, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "qwen-2.5-72b", Name: "Qwen2.5 72B", HF: "Qwen/Qwen2.5-72B-Instruct", Family: "Qwen", ParamsB: 72.7, ActiveB: 72.7, Layers: 80, Hidden: 8192, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "mistral-large", Name: "Mistral Large", HF: "mistralai/Mistral-Large-Instruct-2411", Family: "Mistral", ParamsB: 123, ActiveB: 123, Layers: 88, Hidden: 12288, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
	{ID: "mixtral-8x22b", Name: "Mixtral 8x22B", HF: "mistralai/Mixtral-8x22B-Instruct-v0.1", Family: "Mistral", ParamsB: 141, ActiveB: 39, Layers: 56, Hidden: 6144, KVHeads: 8, HeadDim: 128, MoE: true, Attn: AttnGQA},
	{ID: "deepseek-v3", Name: "DeepSeek-V3", HF: "deepseek-ai/DeepSeek-V3", Family: "DeepSeek", ParamsB: 671, ActiveB: 37, Layers: 61, Hidden: 7168, KVHeads: 1, HeadDim: 128, MoE: true, Attn: AttnMLA, MLALatent: 512, MLARopeDim: 64},
}

var GPURegistry = []GPUSpec{
	{ID: "L40S", Name: "NVIDIA L40S 48GB", VramGB: 48, BandGBs: 864, TFLOPS16: 366},
	{ID: "A100", Name: "NVIDIA A100 80GB", VramGB: 80, BandGBs: 2039, TFLOPS16: 312},
	{ID: "H100", Name: "NVIDIA H100 80GB", VramGB: 80, BandGBs: 3350, TFLOPS16: 990},
	{ID: "H200", Name: "NVIDIA H200 141GB", VramGB: 141, BandGBs: 4800, TFLOPS16: 990},
	{ID: "B200", Name: "NVIDIA B200 192GB", VramGB: 192, BandGBs: 8000, TFLOPS16: 2250},
}

func ResolveModel(name string) ModelSpec {
	n := strings.ToLower(strings.TrimSpace(name))
	for _, m := range ModelRegistry {
		if n == m.ID || n == strings.ToLower(m.HF) || n == strings.ToLower(m.Name) {
			return m
		}
	}
	switch {
	case strings.Contains(n, "deepseek"):
		return mustModel("deepseek-v3")
	case strings.Contains(n, "405"):
		return mustModel("llama-3.1-405b")
	case strings.Contains(n, "mixtral") || strings.Contains(n, "8x22"):
		return mustModel("mixtral-8x22b")
	case strings.Contains(n, "qwen") && strings.Contains(n, "32"):
		return mustModel("qwen-2.5-32b")
	case strings.Contains(n, "qwen") && strings.Contains(n, "7b"):
		return mustModel("qwen-2.5-7b")
	case strings.Contains(n, "qwen"):
		return mustModel("qwen-2.5-72b")
	case strings.Contains(n, "mistral"):
		return mustModel("mistral-large")
	case strings.Contains(n, "8b"):
		return mustModel("llama-3.1-8b")
	case strings.Contains(n, "70") || strings.Contains(n, "3.3"):
		return mustModel("llama-3.3-70b")
	default:
		s := mustModel("llama-3.3-70b")
		s.ID = "approx-70b"
		s.Name = "Unknown · approximated as 70B GQA"
		return s
	}
}

func ResolveGPU(id string, fallbackVRAM float64) GPUSpec {
	u := strings.ToUpper(strings.TrimSpace(id))
	for _, g := range GPURegistry {
		if g.ID == u {
			return g
		}
	}
	if fallbackVRAM <= 0 {
		fallbackVRAM = 80
	}
	return GPUSpec{ID: "CUSTOM", Name: "Custom GPU", VramGB: fallbackVRAM, BandGBs: 2000, TFLOPS16: 300}
}

func mustModel(id string) ModelSpec {
	for _, m := range ModelRegistry {
		if m.ID == id {
			return m
		}
	}
	return ModelRegistry[1]
}

func KVBytesPerToken(spec ModelSpec, kvBytes int) float64 {
	if kvBytes < 1 {
		kvBytes = 2
	}
	if spec.Attn == AttnMLA {
		latent := spec.MLALatent
		if latent == 0 {
			latent = 512
		}
		rope := spec.MLARopeDim
		if rope == 0 {
			rope = 64
		}
		return float64(spec.Layers * (latent + rope) * kvBytes)
	}
	return 2.0 * float64(spec.Layers*spec.KVHeads*spec.HeadDim*kvBytes)
}

func pickParallelism(numGPUs, layers int, weightsGB, kvGB, gpuGB, util float64) (tp, pp int) {
	if numGPUs < 1 {
		numGPUs = 1
	}
	usable := gpuGB * util
	for tp = numGPUs; tp >= 1; tp /= 2 {
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
		if tp == 1 {
			break
		}
	}
	return numGPUs, 1
}

func CloudRate(provider, gpu string) float64 {
	rates := map[string]map[string]float64{
		"lambda": {"H100": 2.49, "H200": 3.79, "A100": 1.29, "L40S": 0.79, "B200": 4.99},
		"runpod": {"H100": 2.29, "H200": 3.49, "A100": 1.19, "L40S": 0.74, "B200": 4.49},
		"aws":    {"H100": 4.10, "H200": 5.40, "A100": 3.67, "L40S": 1.80, "B200": 6.80},
		"gcp":    {"H100": 3.99, "H200": 5.20, "A100": 3.50, "L40S": 1.75, "B200": 6.50},
	}
	p := strings.ToLower(provider)
	g := strings.ToUpper(gpu)
	if table, ok := rates[p]; ok {
		if rate, ok := table[g]; ok {
			return rate
		}
	}
	return 2.49
}

func Normalize(cfg ServingConfig) ServingConfig {
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
	if cfg.DtypeBytes != 1 && cfg.DtypeBytes != 2 {
		cfg.DtypeBytes = 2
	}
	if cfg.KVDtypeBytes != 1 && cfg.KVDtypeBytes != 2 {
		cfg.KVDtypeBytes = cfg.DtypeBytes
	}
	if cfg.Engine == "" {
		cfg.Engine = "vllm"
	}
	if cfg.Provider == "" {
		cfg.Provider = "lambda"
	}
	if cfg.GpuType == "" {
		cfg.GpuType = "H100"
	}
	if cfg.PrefixHit < 0 {
		cfg.PrefixHit = 0
	}
	if cfg.PrefixHit > 0.95 {
		cfg.PrefixHit = 0.95
	}
	return cfg
}

func runtimeGB(engine string) float64 {
	if strings.ToLower(engine) == "sglang" {
		return 2.4
	}
	return 2.2
}

func Analyze(cfg ServingConfig) Analysis {
	cfg = Normalize(cfg)
	spec := ResolveModel(cfg.ModelName)
	gpu := ResolveGPU(cfg.GpuType, cfg.GpuVramGB)
	cfg.GpuVramGB = gpu.VramGB
	weightGB := spec.ParamsB * float64(cfg.DtypeBytes)
	perTok := KVBytesPerToken(spec, cfg.KVDtypeBytes)
	rawKV := perTok * float64(cfg.MaxModelLen) * float64(cfg.MaxNumSeqs)
	eff := 1.0 - cfg.PrefixHit*0.75
	kvGB := (rawKV * eff) / (1024 * 1024 * 1024)
	rt := runtimeGB(cfg.Engine)
	tp, pp := pickParallelism(cfg.NumGpus, spec.Layers, weightGB, kvGB, gpu.VramGB, cfg.GpuMemoryUtilization)
	perGPU := (weightGB / float64(tp*pp)) + (kvGB / float64(tp)) + rt
	usable := gpu.VramGB * cfg.GpuMemoryUtilization
	oom := perGPU > usable
	total := weightGB + kvGB + rt
	kvShare := 0.0
	if total > 0 {
		kvShare = kvGB / total
	}
	maxSeqs := maxFitSeqs(spec, cfg, gpu, weightGB, rt)
	hourly := CloudRate(cfg.Provider, gpu.ID) * float64(cfg.NumGpus)
	a := Analysis{Spec: spec, GPU: gpu, WeightGB: round1(weightGB), KVGB: round1(kvGB), KVPerTokenKB: round1(perTok / 1024), RuntimeGB: rt, TotalGB: round1(total), PerGPUGB: round1(perGPU), UsablePerGPU: round1(usable), TP: tp, PP: pp, OOM: oom, HeadroomGB: round1(usable - perGPU), KVShare: math.Round(kvShare*1000) / 1000, HourlyUSD: math.Round(hourly*100) / 100, MonthlyUSD: math.Round(hourly*24*30*100) / 100, MaxFitSeqs: maxSeqs, PrefixEffective: math.Round(eff*1000) / 1000}
	switch {
	case oom && maxSeqs >= 1:
		a.Recommendation = fmt.Sprintf("OOM at max-num-seqs=%d. Same box fits about %d concurrent sequences if you cut the batch, not the model.", cfg.MaxNumSeqs, maxSeqs)
		a.Fix = fmt.Sprintf("Set --max-num-seqs %d or shorten max-model-len.", maxSeqs)
	case oom:
		a.Recommendation = fmt.Sprintf("OOM: %.1f GB/GPU > %.1f GB usable. Add GPUs or drop to FP8.", perGPU, usable)
		a.Fix = "Add GPUs or switch precision to FP8."
	case a.KVGB > a.WeightGB*0.5:
		a.Recommendation = fmt.Sprintf("KV is %.1f GB versus %.1f GB weights. The invoice is context x concurrency.", a.KVGB, a.WeightGB)
		a.Fix = "Cap max-model-len and turn on prefix caching."
	case spec.MoE && spec.Attn == AttnMLA:
		a.Recommendation = fmt.Sprintf("MoE+MLA: loaded %.0fB, active %.0fB, KV ~ %.1f KB/token.", spec.ParamsB, spec.ActiveB, a.KVPerTokenKB)
		a.Fix = "Do not size KV as if this were GQA."
	case spec.MoE:
		a.Recommendation = fmt.Sprintf("MoE: loaded %.0fB, active %.0fB. Cheap tokens, not free context.", spec.ParamsB, spec.ActiveB)
		a.Fix = "Do not open native context because unit price dropped."
	default:
		a.Recommendation = fmt.Sprintf("TP=%d PP=%d, headroom %.1f GB/GPU. Cluster holds ~%d sequences.", tp, pp, a.HeadroomGB, maxSeqs)
		a.Fix = "Diff engine flags before you pull latest."
	}
	a.VLLMCommand = BuildVLLM(spec.HF, cfg, a)
	a.SGLangCommand = BuildSGLang(spec.HF, cfg, a)
	a.KubeSnippet = BuildKube(spec.HF, cfg, a)
	if strings.ToLower(cfg.Engine) == "sglang" {
		a.EngineCommand = a.SGLangCommand
	} else {
		a.EngineCommand = a.VLLMCommand
	}
	return a
}

func maxFitSeqs(spec ModelSpec, cfg ServingConfig, gpu GPUSpec, weightGB, rt float64) int {
	lo, hi, best := 0, 4096, 0
	for lo <= hi {
		mid := (lo + hi) / 2
		probe := cfg
		probe.MaxNumSeqs = mid
		if mid == 0 {
			lo = mid + 1
			continue
		}
		perTok := KVBytesPerToken(spec, probe.KVDtypeBytes)
		eff := 1.0 - probe.PrefixHit*0.75
		kvGB := (perTok * float64(probe.MaxModelLen) * float64(mid) * eff) / (1024 * 1024 * 1024)
		tp, pp := pickParallelism(probe.NumGpus, spec.Layers, weightGB, kvGB, gpu.VramGB, probe.GpuMemoryUtilization)
		per := (weightGB / float64(tp*pp)) + (kvGB / float64(tp)) + rt
		if per <= gpu.VramGB*probe.GpuMemoryUtilization {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return best
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
