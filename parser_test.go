package main

import (
	"testing"
)

func TestCalculateVRAMAndCost(t *testing.T) {
	cfg := ServingConfig{
		ModelName:            "meta-llama/Llama-3-70b-Instruct",
		MaxModelLen:          32768,
		GpuMemoryUtilization: 0.90,
		NumGpus:              4,
	}

	vram, recommendation := CalculateVRAMAndCost(cfg)

	if vram <= 0 {
		t.Errorf("VRAM hesabı hatalı: %.2f GB", vram)
	}

	if recommendation == "" {
		t.Errorf("Cluster tavsiyesi boş döndü.")
	}
}
