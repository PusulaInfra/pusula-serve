package main

import "fmt"

type ServingConfig struct {
	ModelName            string
	MaxModelLen          int
	GpuMemoryUtilization float64
	MaxNumSeqs           int
}

func AnalyzeConfig(cfg ServingConfig) {
	fmt.Printf("Analiz Edilen Model: %s\n", cfg.ModelName)
	fmt.Printf("Max Model Uzunlugu: %d token\n", cfg.MaxModelLen)
	fmt.Printf("GPU Bellek Kullanimi: %.2f\n", cfg.GpuMemoryUtilization)
	
	if cfg.MaxModelLen > 32768 {
		fmt.Println("UYARI: Yuksek context window KV Cache maliyetini patlatabilir!")
	} else {
		fmt.Println("Konfigurasyon guvenli sinirlarda.")
	}
}
