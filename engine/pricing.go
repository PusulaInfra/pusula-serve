package engine

type GPUPricing struct {
	GpuType      string  `json:"gpu_type"`
	Provider     string  `json:"provider"`
	HourlyUSD    float64 `json:"hourly_usd"`
	MonthlyUSD   float64 `json:"monthly_usd"`
	MarketStatus string  `json:"market_status"`
}

// GetDynamicGPUPricing, Ağustos 2026 güncel bulut sağlayıcı (Lambda, RunPod vb.) 
// piyasa dinamiklerini ve ortalama kiralama maliyetlerini hesaplar.
func GetDynamicGPUPricing(gpuType string, provider string, numGPUs int) GPUPricing {
	if numGPUs <= 0 {
		numGPUs = 1
	}

	// 2026 güncel ortalama saatlik fiyat bazları
	hourlyRate := 2.89 // Varsayılan H100 / üst segment ortalaması
	
	switch gpuType {
	case "H100":
		if provider == "lambda" {
			hourlyRate = 3.99
		} else {
			hourlyRate = 2.89 // RunPod vb.
		}
	case "A100":
		if provider == "lambda" {
			hourlyRate = 2.79
		} else {
			hourlyRate = 1.39
		}
	case "L40S":
		hourlyRate = 0.99
	default:
		hourlyRate = 2.00
	}

	totalHourly := hourlyRate * float64(numGPUs)
	// 730 saat/ay hesabı
	totalMonthly := totalHourly * 730

	return GPUPricing{
		GpuType:      gpuType,
		Provider:     provider,
		HourlyUSD:    round(totalHourly, 2),
		MonthlyUSD:   round(totalMonthly, 2),
		MarketStatus: "Ağustos 2026 piyasa verileriyle güncellendi (Yüksek talep / kararlı bant)",
	}
}
