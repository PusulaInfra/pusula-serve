package engine

type OverheadConfig struct {
	HiddenSize int `json:"hidden_size"`
	NLayers    int `json:"n_layers"`
}

// CalculateOverheadMemory, vLLM ve SGLang altyapılarında model ağırlıkları ve KV cache 
// haricinde tüketilen Activation, Temporary Workspace Buffer ve CUDA Graph VRAM maliyetini hesaplar.
func CalculateOverheadMemory(cfg OverheadConfig, concurrency int, tensorParallelSize int) float64 {
	if tensorParallelSize < 1 {
		tensorParallelSize = 1
	}

	hiddenSize := cfg.HiddenSize
	if hiddenSize <= 0 {
		hiddenSize = 4096 // Varsayılan tipik hidden size
	}

	// 1. Activation ve Temporary Buffer maliyeti (Concurrency ve hidden size ile ölçeklenir)
	// Byte cinsinden geçici çalışma alanı
	activationBytes := float64(concurrency) * float64(hiddenSize) * 2.0 * 4.0
	activationMB := activationBytes / (1024.0 * 1024.0)

	// 2. CUDA Graph rezervasyonu (Genellikle vLLM/SGLang standart modda ~1GB ila 2GB arası sabit VRAM tutar)
	cudaGraphMB := 1024.0 

	totalOverheadMB := activationMB + cudaGraphMB
	totalOverheadGB := (totalOverheadMB / 1024.0) / float64(tensorParallelSize)

	return round(totalOverheadGB, 2)
}
