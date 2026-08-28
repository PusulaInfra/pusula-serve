package serve

import "testing"

func TestLlama70BShortContextFits4x80(t *testing.T) {
	a := Analyze(ServingConfig{
		ModelName:            "meta-llama/Llama-3.3-70B-Instruct",
		MaxModelLen:          8192,
		GpuMemoryUtilization: 0.90,
		NumGpus:              4,
		MaxNumSeqs:           8,
		DtypeBytes:           2,
		GpuType:              "H100",
	})
	if a.WeightGB < 130 || a.WeightGB > 150 {
		t.Fatalf("70B fp16 weights expected ~141GB, got %.1f", a.WeightGB)
	}
	if a.OOM {
		t.Fatalf("8K / 4x80 should fit, perGPU=%.1f", a.PerGPUGB)
	}
}

func TestLongContextInflatesKV(t *testing.T) {
	short := Analyze(ServingConfig{ModelName: "llama-3.3-70b", MaxModelLen: 8192, NumGpus: 4, MaxNumSeqs: 16, GpuMemoryUtilization: 0.9, DtypeBytes: 2, GpuType: "H100"})
	long := Analyze(ServingConfig{ModelName: "llama-3.3-70b", MaxModelLen: 131072, NumGpus: 4, MaxNumSeqs: 16, GpuMemoryUtilization: 0.9, DtypeBytes: 2, GpuType: "H100"})
	if long.KVGB <= short.KVGB*4 {
		t.Fatalf("128K KV should dwarf 8K: short=%.1f long=%.1f", short.KVGB, long.KVGB)
	}
}

func TestDeepSeekIsMoEAndMLA(t *testing.T) {
	a := Analyze(ServingConfig{ModelName: "deepseek-ai/DeepSeek-V3", MaxModelLen: 8192, NumGpus: 8, DtypeBytes: 2, GpuMemoryUtilization: 0.9, MaxNumSeqs: 8, GpuType: "H100"})
	if !a.Spec.MoE || a.Spec.ParamsB < 600 || a.Spec.Attn != AttnMLA {
		t.Fatalf("DeepSeek-V3 must be MoE+MLA, got %+v", a.Spec)
	}
	gqa := Analyze(ServingConfig{ModelName: "llama-3.3-70b", MaxModelLen: 8192, NumGpus: 8, DtypeBytes: 2, MaxNumSeqs: 8, GpuType: "H100"})
	if a.KVPerTokenKB >= gqa.KVPerTokenKB {
		t.Fatalf("MLA KV/token should be smaller than 70B GQA: mla=%.1f gqa=%.1f", a.KVPerTokenKB, gqa.KVPerTokenKB)
	}
}

func Test405BIsNotApprox70B(t *testing.T) {
	a := Analyze(ServingConfig{ModelName: "llama-3.1-405b", NumGpus: 8, GpuType: "H100"})
	if a.Spec.ParamsB < 400 {
		t.Fatalf("405B should not collapse to 70B, got %+v", a.Spec)
	}
}
