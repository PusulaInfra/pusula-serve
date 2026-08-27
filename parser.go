package main

import "fmt"

type ServingConfig struct {
	ModelName            string
	MaxModelLen          int
	GpuMemoryUtilization float64
	NumGpus              int
}

// Akilli Bellek ve KV Cache Maliyet Hesaplayici
func CalculateVRAMAndCost(cfg ServingConfig) (float64, string) {
	// Basit bir kestirim: 70B model yaklasik 140GB (FP16) yer kaplar. 
	// KV Cache uzunluga ve paralellige gore artar.
	baseModelGB := 140.0 
	if len(cfg.ModelName) < 20 {
		baseModelGB = 30.0 // Daha kucuk modeller icin
	}
	
	kvCacheGB := float64(cfg.MaxModelLen) * 0.0005 * float64(cfg.NumGpus)
	totalEstimatedVRAM := baseModelGB + kvCacheGB
	
	recommendation := "A100 / H100 80GB GPUCluster onerilir."
	if totalEstimatedVRAM > 160 {
		recommendation = "UYARI: Yuksek context nedeniyle en az 4x A100 (80GB) veya Multi-Node GPU gereklidir!"
	}

	return totalEstimatedVRAM, recommendation
}

func AnalyzeConfig(cfg ServingConfig) {
	vram, rec := CalculateVRAMAndCost(cfg)
	fmt.Printf("Model: %s\n", cfg.ModelName)
	fmt.Printf("Tahmini VRAM Kullanimi: %.2f GB\n", vram)
	fmt.Println("Cluster Onerisi:", rec)
}
