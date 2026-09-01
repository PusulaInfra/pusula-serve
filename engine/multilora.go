package engine

type LoRAConfig struct {
	NumActiveAdapters int `json:"num_active_adapters"`
	Rank              int `json:"rank"`
	HiddenSize        int `json:"hidden_size"`
	NumLayers         int `json:"num_layers"`
	BytesPerParam     int `json:"bytes_per_param"` // 2 = FP16/BF16
}

type LoRAMemoryResult struct {
	ExtraVRAMGB float64 `json:"extra_vram_gb"`
	Description string  `json:"description"`
}

// CalculateLoRAMemory, ana modelin üzerine dinamik yüklenen LoRA adaptörlerinin 
// VRAM üzerindeki ek bellek maliyetini hesaplar.
func CalculateLoRAMemory(cfg LoRAConfig) LoRAMemoryResult {
	if cfg.Rank <= 0 { cfg.Rank = 16 }
	if cfg.HiddenSize <= 0 { cfg.HiddenSize = 4096 }
	if cfg.NumLayers <= 0 { cfg.NumLayers = 32 }
	if cfg.BytesPerParam <= 0 { cfg.BytesPerParam = 2 }

	// LoRA parametre formülü: 2 * (hidden_size * rank) * num_layers (her adaptör için)
	paramsPerAdapter := 2 * int64(cfg.HiddenSize) * int64(cfg.Rank) * int64(cfg.NumLayers)
	totalBytesPerAdapter := paramsPerAdapter * int64(cfg.BytesPerParam)
	totalVRAMBytes := totalBytesPerAdapter * int64(cfg.NumActiveAdapters)
	extraGB := float64(totalVRAMBytes) / (1024 * 1024 * 1024)

	return LoRAMemoryResult{
		ExtraVRAMGB: round(extraGB, 2),
		Description: "Multi-LoRA eşzamanlı adaptör yükü için VRAM ayrımı yapıldı.",
	}
}
