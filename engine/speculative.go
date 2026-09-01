package engine

type SpeculativeSimulation struct {
	DraftModelName  string
	TargetModelName string
	AcceptanceRate  float64 // Ortalama kabul oranı (örn: 0.70)
	SpeedupFactor   float64
}

func SimulateSpeculativeDecoding(acceptanceRate float64) SpeculativeSimulation {
	// Temel formül: 1 / (1 - acceptanceRate * draft_ratio) tahmini çarpan
	factor := 1.0 / (1.0 - (acceptanceRate * 0.8))
	return SpeculativeSimulation{
		DraftModelName:  "Llama-3-8B-Instruct",
		TargetModelName: "Llama-3.1-70B-Instruct",
		AcceptanceRate:  acceptanceRate,
		SpeedupFactor:   factor,
	}
}
