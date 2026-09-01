package engine

import (
	"fmt"
)

type MoEConfig struct {
	IsMoE         bool    `json:"is_moe"`
	TotalParamsB  float64 `json:"total_params_b"`
	ActiveParamsB float64 `json:"active_params_b"`
	NumExperts    int     `json:"num_experts"`
}

type MoEResult struct {
	IsMoE                 bool   `json:"is_moe"`
	RecommendedEPSize     int    `json:"recommended_ep_size"`
	CommunicationOverhead string `json:"communication_overhead"`
	MemoryLoadNote        string `json:"memory_load_note"`
}

// CalculateMoEAndEP, DeepSeek-V3, Mixtral ve Qwen3-MoE gibi Mixture of Experts 
// modelleri için Loaded vs Active bellek dağılımını ve en uygun Expert Parallelism (EP) boyutunu hesaplar.
func CalculateMoEAndEP(cfg MoEConfig, totalGpus int) MoEResult {
	if !cfg.IsMoE {
		return MoEResult{
			IsMoE:                 false,
			RecommendedEPSize:     1,
			CommunicationOverhead: "Yok",
			MemoryLoadNote:        "Dense model, EP gerekmez.",
		}
	}

	// Uzman paralelliği (EP) toplam GPU'ya eşit/küçük ve uzman sayısının bir böleni olmalıdır.
	var possibleEPs []int
	candidateEPs := []int{1, 2, 4, 8, 16, 32}
	for _, ep := range candidateEPs {
		if ep <= totalGpus && cfg.NumExperts > 0 && cfg.NumExperts%ep == 0 {
			possibleEPs = append(possibleEPs, ep)
		}
	}

	optimalEP := 1
	if len(possibleEPs) > 0 {
		optimalEP = possibleEPs[len(possibleEPs)-1]
	}

	commOverhead := "Düşük / Orta"
	if optimalEP > 4 {
		commOverhead = "Yüksek (NVLink / InfiniBand önerilir)"
	}

	memoryNote := fmt.Sprintf("Model belleğe tam yüklenir (%.1fB), token işleme anında ise %.1fB aktif olur.", cfg.TotalParamsB, cfg.ActiveParamsB)

	return MoEResult{
		IsMoE:                 true,
		RecommendedEPSize:     optimalEP,
		CommunicationOverhead: commOverhead,
		MemoryLoadNote:        memoryNote,
	}
}
