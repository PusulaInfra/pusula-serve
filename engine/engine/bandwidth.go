package engine

type HBMBoundResult struct {
	ModelSizeGB          float64 `json:"model_size_gb"`
	GPUBankwidthTBS      float64 `json:"gpu_bandwidth_tb_s"`
	MaxTheoreticalTokSec float64 `json:"max_theoretical_tok_sec"`
	IsHBMBound           bool    `json:"is_hbm_bound"`
	Explanation          string  `json:"explanation"`
}

// CalculateHBMBoundTokSec, GPU'nun bellek bant genişliği ile model boyutu arasındaki ilişkiyi 
// kullanarak teorik maksimum token/saniye tavanını ve HBM darboğaz durumunu hesaplar.
func CalculateHBMBoundTokSec(modelSizeGB float64, gpuMemoryBandwidthTBS float64) HBMBoundResult {
	if gpuMemoryBandwidthTBS <= 0 {
		gpuMemoryBandwidthTBS = 3.35 // Örn: H100 SXM tipik bant genişliği (TB/s)
	}
	if modelSizeGB <= 0 {
		modelSizeGB = 1.0
	}

	// Bandwidth TB/s -> GB/s çevrimi (1 TB = 1024 GB)
	bwGBS := gpuMemoryBandwidthTBS * 1024.0

	// Decode sırasında her adımda tüm model ağırlıkları HBM'den okunur:
	// Max Tok/s = Bellek Bant Genişliği (GB/s) / Model Boyutu (GB)
	maxTokSec := bwGBS / modelSizeGB

	return HBMBoundResult{
		ModelSizeGB:          round(modelSizeGB, 2),
		GPUBankwidthTBS:      gpuMemoryBandwidthTBS,
		MaxTheoreticalTokSec: round(maxTokSec, 1),
		IsHBMBound:           true,
		Explanation:          "Decode aşaması tamamen HBM (Memory Bandwidth) sınırlarındadır; teorik tavan aşlamaz.",
	}
}
