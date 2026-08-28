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
