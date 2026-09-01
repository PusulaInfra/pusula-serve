package engine

type WorkloadPrefixProfile struct {
	WorkloadType  string  `json:"workload_type"`
	DefaultHitRate float64 `json:"default_hit_rate"`
	Description   string  `json:"description"`
}

var WorkloadProfiles = map[string]WorkloadPrefixProfile{
	"Agent / Multi-turn Chat": {
		WorkloadType:   "Agent / Multi-turn Chat",
		DefaultHitRate: 0.65,
		Description:    "Sistem promptları, araç tanımları ve uzun konuşma geçmişi nedeniyle yüksek önbellek isabeti (%65).",
	},
	"RAG / Document Q&A": {
		WorkloadType:   "RAG / Document Q&A",
		DefaultHitRate: 0.40,
		Description:    "Her sorguda farklı doküman parçaları bağlama eklendiği için orta düzey isabet (%40).",
	},
	"Raw Completion / Batch": {
		WorkloadType:   "Raw Completion / Batch",
		DefaultHitRate: 0.05,
		Description:    "Tekil ve bağımsız istekler olduğu için önbellek tekrarı minimum düzeyde (%5).",
	},
}

// GetWorkloadPrefixDefault, agent, RAG veya toplu işlem (batch) gibi 
// farklı iş yüklerine göre en akıllı prefix-hit (önbellek isabet) varsayılanlarını döner.
func GetWorkloadPrefixDefault(workloadType string) WorkloadPrefixProfile {
	profile, exists := WorkloadProfiles[workloadType]
	if !exists {
		// Varsayılan genel agent eğilimi
		return WorkloadProfiles["Agent / Multi-turn Chat"]
	}
	return profile
}
