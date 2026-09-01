package engine

type GPUFleetSpec struct {
	FleetName        string  `json:"fleet_name"`
	VramPerGPU       float64 `json:"vram_per_gpu"`
	BandwidthTBS     float64 `json:"bandwidth_tb_s"`
	RelativeCompute  float64 `json:"relative_compute"`
	DefaultGPUCount  int     `json:"default_gpu_count"`
}

var AvailableFleets = map[string]GPUFleetSpec{
	"4x H100 80GB":     {FleetName: "4x H100 80GB", VramPerGPU: 80.0, BandwidthTBS: 3.35, RelativeCompute: 10.0, DefaultGPUCount: 4},
	"8x RTX 5090 32GB": {FleetName: "8x RTX 5090 32GB", VramPerGPU: 32.0, BandwidthTBS: 1.79, RelativeCompute: 7.5, DefaultGPUCount: 8},
	"2x H200 141GB":    {FleetName: "2x H200 141GB", VramPerGPU: 141.0, BandwidthTBS: 4.80, RelativeCompute: 6.0, DefaultGPUCount: 2},
	"8x H100 80GB":     {FleetName: "8x H100 80GB", VramPerGPU: 80.0, BandwidthTBS: 3.35, RelativeCompute: 20.0, DefaultGPUCount: 8},
}

type FleetComparisonResult struct {
	FleetName   string  `json:"fleet_name"`
	CanFit      bool    `json:"can_fit"`
	TotalVRAMGB float64 `json:"total_vram_gb"`
	BandwidthTBS float64 `json:"bandwidth_tb_s"`
	StatusNote  string  `json:"status_note"`
}

// CompareFleets, farklı GPU kümelerinin (H100, RTX 5090, H200 vb.) 
// toplam VRAM ve bant genişliği açısından modele uygunluğunu çoklu olarak kıyaslar.
func CompareFleets(totalModelAndKVSizeGB float64) []FleetComparisonResult {
	var results []FleetComparisonResult

	for _, spec := range AvailableFleets {
		totalClusterVRAM := spec.VramPerGPU * float64(spec.DefaultGPUCount)
		canFit := totalClusterVRAM >= totalModelAndKVSizeGB

		note := "Uygun (Konfigürasyon çalışır)"
		if !canFit {
			note = "Yetersiz VRAM (OOM riski var)"
		}

		results = append(results, FleetComparisonResult{
			FleetName:    spec.FleetName,
			CanFit:       canFit,
			TotalVRAMGB:  totalClusterVRAM,
			BandwidthTBS: spec.BandwidthTBS,
			StatusNote:   note,
		})
	}

	return results
}
