package engine

import (
	"net/url"
	"strconv"
)

type AppState struct {
	ModelName   string `json:"model_name"`
	Concurrency int    `json:"concurrency"`
	MaxContext  int    `json:"max_context"`
	QuantType   string `json:"quant_type"`
	TotalGPUs   int    `json:"total_gpus"`
	Workload    string `json:"workload"`
}

// EncodeStateToURLParams, hesaplama parametrelerini URL query string formatına çevirir.
func EncodeStateToURLParams(state AppState) string {
	params := url.Values{}
	params.Set("model", state.ModelName)
	params.Set("conc", strconv.Itoa(state.Concurrency))
	params.Set("ctx", strconv.Itoa(state.MaxContext))
	params.Set("quant", state.QuantType)
	params.Set("gpus", strconv.Itoa(state.TotalGPUs))
	params.Set("workload", state.Workload)
	return params.Encode()
}

// DecodeStateFromURLParams, URL'den gelen parametreleri okuyarak uygulama state'ine dönüştürür.
func DecodeStateFromURLParams(queryValues url.Values) AppState {
	conc, _ := strconv.Atoi(queryValues.Get("conc"))
	if conc <= 0 {
		conc = 32
	}

	ctx, _ := strconv.Atoi(queryValues.Get("ctx"))
	if ctx <= 0 {
		ctx = 4096
	}

	gpus, _ := strconv.Atoi(queryValues.Get("gpus"))
	if gpus <= 0 {
		gpus = 4
	}

	model := queryValues.Get("model")
	if model == "" {
		model = "Llama-3.1-70B"
	}

	quant := queryValues.Get("quant")
	if quant == "" {
		quant = "FP16"
	}

	workload := queryValues.Get("workload")
	if workload == "" {
		workload = "Agent / Multi-turn Chat"
	}

	return AppState{
		ModelName:   model,
		Concurrency: conc,
		MaxContext:  ctx,
		QuantType:   quant,
		TotalGPUs:   gpus,
		Workload:    workload,
	}
}
