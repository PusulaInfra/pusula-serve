package main

func main() {
	cfg := ServingConfig{
		ModelName:            "meta-llama/Llama-3-70b-Instruct",
		MaxModelLen:          65536,
		GpuMemoryUtilization: 0.90,
		MaxNumSeqs:           256,
	}

	AnalyzeConfig(cfg)
}
