package engine

type PDSplitRecommendation struct {
	PDSplitRecommended bool   `json:"pd_split_recommended"`
	Strategy           string `json:"strategy"`
	Reason             string `json:"reason"`
}

// RecommendPDSplit, kaba "yes/maybe/no" yaklaşımı yerine concurrency, 
// ortalama prompt uzunluğu ve generation uzunluğu oranına bakarak dinamik Prefill/Decode split önerisi yapar.
func RecommendPDSplit(concurrency int, avgPromptLen int, avgGenLen int) PDSplitRecommendation {
	if avgGenLen < 1 {
		avgGenLen = 1
	}

	ratio := float64(avgPromptLen) / float64(avgGenLen)

	// Yüksek concurrency ve ağır prompt oranı Prefill / Decode düğümlerinin ayrılmasını (Disaggregation) zorunlu kılar
	if concurrency > 64 && ratio > 3.0 {
		return PDSplitRecommendation{
			PDSplitRecommended: true,
			Strategy:           "Disaggregated P/D (Prefill ve Decode havuzları ayrılmalı)",
			Reason:             "Yüksek eşzamanlılık ve uzun promptlar, decode aşamasında prefill gecikmelerine (head-of-line blocking) yol açıyor.",
		}
	} else if concurrency > 32 && ratio > 1.5 {
		return PDSplitRecommendation{
			PDSplitRecommended: true,
			Strategy:           "Hibrit / Hafif P/D Split Önerilir",
			Reason:             "İş yükü yoğunluğu artmaya başladı; uzun bağlamlar için önbellek izolasyonu faydalı olacaktır.",
		}
	}

	return PDSplitRecommendation{
				PDSplitRecommended: false,
		Strategy:           "Co-located (Standart tekil küme)",
		Reason:             "İş yükü dengeli; ek mimari karmaşıklığa ve ağ maliyetine gerek yok.",
	}
}
