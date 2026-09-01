package engine

type SLAResult struct {
	TTFTMs float64 `json:"ttft_ms"`
	TPOTMs float64 `json:"tpot_ms"`
	Status string  `json:"status"`
}

// SimulateSLALatency, prompt uzunluğuna ve eşzamanlılığa bağlı olarak 
// TTFT (İlk Token Süresi) ve TPOT (Token Başına Süre) gecikmelerini öngörür.
func SimulateSLALatency(promptLen int, genLen int, concurrency int) SLAResult {
	if promptLen <= 0 { promptLen = 1024 }
	if genLen <= 0 { genLen = 128 }
	if concurrency <= 0 { concurrency = 16 }

	// Prefill aşaması (TTFT) hesaplaması
	ttft := float64(promptLen) * 0.015 * (1.0 + float64(concurrency)*0.02)
	
	// Decode aşaması (TPOT) hesaplaması
	tpot := float64(genLen) * 12.5 * (1.0 + float64(concurrency)*0.01)

	status := "SLA Hedefleri İçinde"
	if ttft > 1000.0 || tpot > 50.0 {
		status = "Yüksek Eşzamanlılık Gecikme Riski"
	}

	return SLAResult{
		TTFTMs: round(ttft, 1),
		TPOTMs: round(tpot, 1),
		Status: status,
	}
}
