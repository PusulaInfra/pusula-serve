package engine

type GreenAIResult struct {
	TotalPowerKW           float64 `json:"total_power_kw"`
	MonthlyKWh             float64 `json:"monthly_kwh"`
	MonthlyPowerUSD        float64 `json:"monthly_power_usd"`
	CarbonEmissionsKgMonth float64 `json:"carbon_emissions_kg_month"`
}

// CalculateGreenAI, GPU kümesinin TDP güç tüketimini, aylık elektrik maliyetini 
// ve karbon salınımını hesaplar.
func CalculateGreenAI(numGPUs int, avgPowerKWPerGPU float64, electricityCostKWh float64) GreenAIResult {
	if avgPowerKWPerGPU <= 0 {
		avgPowerKWPerGPU = 0.7 // Örn: H100 TDP ~ 700W (0.7 kW)
	}
	if electricityCostKWh <= 0 {
		electricityCostKWh = 0.12 // kWh başına ortalama endüstriyel maliyet ($)
	}
	if numGPUs <= 0 { numGPUs = 1 }

	totalPowerKW := avgPowerKWPerGPU * float64(numGPUs)
	monthlyKWh := totalPowerKW * 730.0 // 730 saat/ay
	monthlyPowerUSD := monthlyKWh * electricityCostKWh
	
	// Küresel ortalama şebeke emisyon faktörü: ~0.4 kg CO2 / kWh
	carbonKg := monthlyKWh * 0.4

	return GreenAIResult{
		TotalPowerKW:           round(totalPowerKW, 2),
		MonthlyKWh:             round(monthlyKWh, 1),
		MonthlyPowerUSD:        round(monthlyPowerUSD, 2),
		CarbonEmissionsKgMonth: round(carbonKg, 1),
	}
}
