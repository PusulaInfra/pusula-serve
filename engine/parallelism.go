package engine

type ParallelismOptimizationResult struct {
	TensorParallelSize         int     `json:"tensor_parallel_size"`
	PipelineParallelSize       int     `json:"pipeline_parallel_size"`
	EstimatedPipelineBubblePct float64 `json:"estimated_pipeline_bubble_pct"`
	CommunicationProfile       string  `json:"communication_profile"`
}

// OptimizeParallelism, basit power-of-2 mantığı yerine pipeline bubble maliyetini, 
// TP inter-GPU ağ yükünü (NVLink / Ethernet) ve donanım sınırlarını optimize ederek en ideal paralel stratejiyi seçer.
func OptimizeParallelism(totalGpus int, modelParamsB float64) ParallelismOptimizationResult {
	if totalGpus < 1 {
		totalGpus = 1
	}

	bestTP := 1
	bestPP := 1
	minPenalty := 999999.0

	// Olası Tensor Parallel (TP) ve Pipeline Parallel (PP) kombinasyonları
	tpOptions := []int{1, 2, 4, 8}
	for _, tp := range tpOptions {
		if tp > totalGpus || totalGpus%tp != 0 {
			continue
		}
		pp := totalGpus / tp

		// Pipeline bubble maliyeti (PP arttıkça artar)
		bubblePenalty := float64(pp-1) * 0.08
		
		// TP > 4 olduğunda inter-node iletişim yükü (NVLink eksikliğinde Ethernet darboğazı) artar
		tpPenalty := 0.0
		if tp > 4 {
			tpPenalty = 0.20
		}

		totalPenalty := bubblePenalty + tpPenalty
		if totalPenalty < minPenalty {
			minPenalty = totalPenalty
			bestTP = tp
			bestPP = pp
		}
	}

	// Pipeline Bubble yüzde tahmini
	bubblePct := 0.0
	if bestPP > 1 {
		bubblePct = round((float64(bestPP-1)/32.0)*100.0, 1)
	}

	commProfile := "Standart PCIe / Düşük Ağ Yükü"
	if bestTP > 2 {
		commProfile = "NVLink / InfiniBand Şart (Yüksek bant genişliği gerekli)"
	}

	return ParallelismOptimizationResult{
		TensorParallelSize:         bestTP,
		PipelineParallelSize:       bestPP,
		EstimatedPipelineBubblePct: bubblePct,
		CommunicationProfile:       commProfile,
	}
}
