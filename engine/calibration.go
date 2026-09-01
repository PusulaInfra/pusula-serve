package engine

// CalibrateKVAndRuntime, vLLM ve SGLang ölçümlerine dayanarak kaba weight*0.08 
// yaklaşımını gerçekçi KV cache boyutu ve runtime metrikleriyle değiştirir.
func CalibrateKVAndRuntime(
	nLayers int,
	numKvHeads int,
	headDim int,
	concurrency int,
	maxContextLen int,
	tensorParallelSize int,
	kvCacheDtype string,
	baseLatencyMs float64,
) (float64, float64, float64) {
	bytesPerElem := 2.0
	if kvCacheDtype == "fp8" || kvCacheDtype == "int8" {
		bytesPerElem = 1.0
	}

	if tensorParallelSize < 1 {
		tensorParallelSize = 1
	}

	// 1. Katman başına düşen KV boyutu (K ve V için 2 çarpanı)
	kvBytesPerTokenPerLayer := 2.0 * float64(numKvHeads) * float64(headDim) * bytesPerElem
	totalKvPerToken := (kvBytesPerTokenPerLayer * float64(nLayers)) / float64(tensorParallelSize)

	// Toplam KV Cache (PagedAttention %4 fragmentasyon payı dahil)
	totalKvCacheBytes := float64(concurrency) * float64(maxContextLen) * totalKvPerToken * 1.04
	totalKvCacheGB := totalKvCacheBytes / (1024 * 1024 * 1024)

	if baseLatencyMs <= 0 {
		baseLatencyMs = 15.0
	}

	// 2. Yüksek concurrency ve uzun context durumlarındaki non-linear maliyet artışı
	concurrencyPenalty := 1.0 + (float64(concurrency) * 0.002)
	contextPenalty := (float64(maxContextLen) / 2048.0) * 0.35

	estimatedLatencyMs := baseLatencyMs * concurrencyPenalty * (1.0 + contextPenalty)
	
	lat := estimatedLatencyMs
	if lat < 1.0 {
		lat = 1.0
	}
	estimatedTokSecPerGPU := (float64(concurrency) * 1000.0) / lat

	return totalKvCacheGB, estimatedLatencyMs, estimatedTokSecPerGPU
}
