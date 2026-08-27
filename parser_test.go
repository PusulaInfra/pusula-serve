package main

import "testing"

func TestLlama70BShortContextFits4x80(t *testing.T) {
	a := Analyze(ServingConfig{
		ModelName:            "meta-llama/Llama-3-70b-Instruct",
		MaxModelLen:          8192,
		GpuMemoryUtilization: 0.90,
		NumGpus:              4,
		MaxNumSeqs:           8,
		DtypeBytes:           2,
	})
	if a.WeightGB < 130 || a.WeightGB > 150 {
		t.Fatalf("70B fp16 weights expected ~141GB, got %.1f", a.WeightGB)
	}
	if a.OOM {
		t.Fatalf("8K / 4x80 should fit, perGPU=%.1f", a.PerGPUGB)
	}
}

func TestLongContextInflatesKV(t *testing.T) {
	short := Analyze(ServingConfig{ModelName: "llama-3-70b", MaxModelLen: 8192, NumGpus: 4, MaxNumSeqs: 16, GpuMemoryUtilization: 0.9, DtypeBytes: 2})
	long := Analyze(ServingConfig{ModelName: "llama-3-70b", MaxModelLen: 131072, NumGpus: 4, MaxNumSeqs: 16, GpuMemoryUtilization: 0.9, DtypeBytes: 2})
	if long.KVGB <= short.KVGB*4 {
		t.Fatalf("128K KV should dwarf 8K: short=%.1f long=%.1f", short.KVGB, long.KVGB)
	}
}

func TestDeepSeekIsMoENot70B(t *testing.T) {
	a := Analyze(ServingConfig{ModelName: "deepseek-ai/DeepSeek-V3", MaxModelLen: 8192, NumGpus: 8, DtypeBytes: 2, GpuMemoryUtilization: 0.9, MaxNumSeqs: 8})
	if !a.Spec.MoE || a.Spec.ParamsB < 600 {
		t.Fatalf("DeepSeek-V3 must use total expert weights, got %+v", a.Spec)
	}
}
