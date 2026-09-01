package engine

// InterconnectType defines GPU communication topology speed.
type InterconnectType string

const (
	NVLink900GBs InterconnectType = "NVLink-900GB/s"
	PCIeGen5     InterconnectType = "PCIe-Gen5-64GB/s"
	StandardPCIe InterconnectType = "PCIe-Gen4-32GB/s"
)

// CalculateKVCacheBytes computes KV memory footprint based on quantization precision.
func CalculateKVCacheBytes(numLayers, numHeads, headDim, ctxLen, maxSeqs int, kvDtypeBytes float64) float64 {
	// KV cache per token = 2 * num_layers * kv_heads * head_dim * dtype_bytes
	bytesPerToken := 2.0 * float64(numLayers) * float64(numHeads) * float64(headDim) * kvDtypeBytes
	totalBytes := bytesPerToken * float64(ctxLen) * float64(maxSeqs)
	return totalBytes / (1024 * 1024 * 1024) // GB cinsinden
}

// CheckInterconnectBottleneck evaluates communication warnings for Tensor Parallelism.
func CheckInterconnectBottleneck(tpSize int, interconnect InterconnectType) string {
	if tpSize > 2 && interconnect == StandardPCIe {
		return "WARNING: High TP size with Standard PCIe will cause significant communication bottleneck!"
	}
	if tpSize > 4 && interconnect == PCIeGen5 {
		return "NOTICE: PCIe Gen5 handles moderate TP, but NVLink is recommended for 70B+ models."
	}
	return "OPTIMAL: Interconnect bandwidth is sufficient for current topology."
}
