package measure

import (
	"github.com/PusulaInfra/pusula-serve/internal/live"
	"github.com/PusulaInfra/pusula-serve/internal/quality"
	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

type Report struct {
	Kind       string          `json:"kind"`
	Disclaimer string          `json:"disclaimer"`
	Estimate   serve.Analysis  `json:"estimate"`
	VRAM       live.Snapshot   `json:"vram"`
	Bench      quality.Result  `json:"bench"`
}

func Run(cfg serve.ServingConfig, benchSec int) Report {
	return Report{
		Kind:       "measure",
		Disclaimer: live.Disclaimer,
		Estimate:   serve.Analyze(cfg),
		VRAM:       live.SnapshotNow(),
		Bench:      quality.Run(benchSec),
	}
}
