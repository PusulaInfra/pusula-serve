package engine

type QuantProfile struct {
	BitWidth            float64
	OverheadMultiplier  float64
	QualityFlag         string
}

var QuantProfiles = map[string]QuantProfile{
	"FP16":        {BitWidth: 16.0, OverheadMultiplier: 1.0, QualityFlag: "Baseline (Kayıpsız)"},
	"FP8":         {BitWidth: 8.0, OverheadMultiplier: 1.05, QualityFlag: "Çok Yüksek (~0 PPL kaybı)"},
	"AWQ-4bit":    {BitWidth: 4.5, OverheadMultiplier: 1.10, QualityFlag: "Yüksek (Minimal kayıp)"},
	"GPTQ-4bit":   {BitWidth: 4.5, OverheadMultiplier: 1.10, QualityFlag: "Yüksek (Minimal kayıp)"},
	"GGUF-Q4_K_M": {BitWidth: 4.8, OverheadMultiplier: 1.02, QualityFlag: "İyi (CPU/GPU hibrit uyumlu)"},
	"GGUF-Q2_K":   {BitWidth: 2.8, OverheadMultiplier: 1.02, QualityFlag: "Düşük (Dikkatli kullanılmalı)"},
}

type QuantResult struct {
	QuantType        string  `json:"quant_type"`
	EstimatedSizeGB  float64 `json:"estimated_size_gb"`
	QualityFlag      string  `json:"quality_flag"`
}

// CalculateQuantizedModelSize, kaba 0.5x yaklaşımını ortadan kaldırarak 
// AWQ, GPTQ, GGUF ve FP8 için gerçekçi boyut ve kalite bayraklarını hesaplar.
func CalculateQuantizedModelSize(paramsB float64, quantType string) QuantResult {
	profile, exists := QuantProfiles[quantType]
	if !exists {
		profile = QuantProfiles["FP16"]
		quantType = "FP16"
	}

	// Boyut formülü: Parametre (Milyar) * (Bit / 16) * overhead * 2 byte (FP16 tabanlı)
	sizeGB := paramsB * (profile.BitWidth / 16.0) * profile.OverheadMultiplier * 2.0

	return QuantResult{
		QuantType:       quantType,
		EstimatedSizeGB: round(sizeGB, 2),
		QualityFlag:     profile.QualityFlag,
	}
}
