package serve

// Registry extensions for GPUs and models mentioned in README but not in code.
// Pusula plans for consumer, datacenter, Apple UMA, and edge hardware.

func initRegistries() {
	// Append extra GPUs
	GPURegistry = append(GPURegistry,
		GPUSpec{ID: "5090", Name: "RTX 5090 32GB", VramGB: 32, BandGBs: 1790, TFLOPS16: 210, Kind: "consumer"},
		GPUSpec{ID: "pro6000", Name: "RTX Pro 6000 96GB", VramGB: 96, BandGBs: 1790, TFLOPS16: 250, Kind: "professional"},
		GPUSpec{ID: "spark", Name: "DGX Spark GB10 128 UMA", VramGB: 128, BandGBs: 273, TFLOPS16: 100, UMA: true, Kind: "edge"},
		GPUSpec{ID: "studio192", Name: "Mac Studio 192 UMA", VramGB: 192, BandGBs: 800, TFLOPS16: 70, UMA: true, Kind: "apple"},
		GPUSpec{ID: "mini16", Name: "Mac Mini 16 UMA", VramGB: 16, BandGBs: 100, TFLOPS16: 8, UMA: true, Kind: "apple"},
		GPUSpec{ID: "jetson", Name: "Jetson AGX 64GB", VramGB: 64, BandGBs: 200, TFLOPS16: 20, UMA: true, Kind: "edge"},
	)

	// Append extra models
	ModelRegistry = append(ModelRegistry,
		ModelSpec{ID: "llama-4-scout", Name: "Llama 4 Scout", HF: "meta-llama/Llama-4-Scout", Family: "Llama", ParamsB: 109, ActiveB: 17, Layers: 48, Hidden: 5120, KVHeads: 8, HeadDim: 128, MoE: true, Attn: AttnGQA},
		ModelSpec{ID: "qwen2.5-vl-72b", Name: "Qwen2.5-VL 72B", HF: "Qwen/Qwen2.5-VL-72B-Instruct", Family: "QwenVL", ParamsB: 73, ActiveB: 73, Layers: 80, Hidden: 8192, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
		ModelSpec{ID: "whisper-large-v3", Name: "Whisper Large v3", HF: "openai/whisper-large-v3", Family: "Whisper", ParamsB: 1.55, ActiveB: 1.55, Layers: 32, Hidden: 1280, KVHeads: 20, HeadDim: 64, Attn: AttnGQA},
		ModelSpec{ID: "e5-mistral-7b", Name: "E5 Mistral 7B embed", HF: "intfloat/e5-mistral-7b-instruct", Family: "Embed", ParamsB: 7.1, ActiveB: 7.1, Layers: 32, Hidden: 4096, KVHeads: 8, HeadDim: 128, Attn: AttnGQA},
		ModelSpec{ID: "bge-reranker-v2", Name: "BGE Reranker v2", HF: "BAAI/bge-reranker-v2-m3", Family: "Rerank", ParamsB: 0.57, ActiveB: 0.57, Layers: 24, Hidden: 1024, KVHeads: 16, HeadDim: 64, Attn: AttnGQA},
	)
}

func init() {
	initRegistries()
}
