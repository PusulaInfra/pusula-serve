package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type APIResponse struct {
	ModelName      string  `json:"model_name"`
	Engine         string  `json:"engine"`
	WeightGB       float64 `json:"weight_gb"`
	KVGB           float64 `json:"kv_gb"`
	VramGB         float64 `json:"vram_gb"`
	PerGPUGB       float64 `json:"per_gpu_gb"`
	TP             int     `json:"tp"`
	PP             int     `json:"pp"`
	OOM            bool    `json:"oom"`
	MonthlyCostUSD float64 `json:"monthly_cost_usd"`
	Command        string  `json:"command"`
	Recommendation string  `json:"recommendation"`
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func buildCommand(engine, model string, cfg ServingConfig, a Analysis) string {
	if engine == "sglang" {
		return fmt.Sprintf(
			"python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --mem-fraction-static %.2f --port 8000",
			model, cfg.MaxModelLen, a.TP, cfg.GpuMemoryUtilization)
	}
	return fmt.Sprintf(
		"vllm serve %s --max-model-len %d --tensor-parallel-size %d --pipeline-parallel-size %d --gpu-memory-utilization %.2f --max-num-seqs %d",
		model, cfg.MaxModelLen, a.TP, a.PP, cfg.GpuMemoryUtilization, cfg.MaxNumSeqs)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	modelName := r.URL.Query().Get("model")
	if modelName == "" {
		modelName = "meta-llama/Llama-3-70b-Instruct"
	}
	engine := r.URL.Query().Get("engine")
	if engine == "" {
		engine = "vllm"
	}
	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          queryInt(r, "max_model_len", 8192),
		NumGpus:              queryInt(r, "num_gpus", 4),
		MaxNumSeqs:           32,
		GpuMemoryUtilization: 0.90,
		DtypeBytes:           2,
	}
	a := Analyze(cfg)
	hourly := 2.15 * float64(cfg.NumGpus)
	json.NewEncoder(w).Encode(APIResponse{
		ModelName:      modelName,
		Engine:         engine,
		WeightGB:       a.WeightGB,
		KVGB:           a.KVGB,
		VramGB:         a.TotalGB,
		PerGPUGB:       a.PerGPUGB,
		TP:             a.TP,
		PP:             a.PP,
		OOM:            a.OOM,
		MonthlyCostUSD: hourly * 24 * 30,
		Command:        buildCommand(engine, modelName, cfg, a),
		Recommendation: a.Recommendation,
	})
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	modelName := "meta-llama/Llama-3-70b-Instruct"
	maxModelLen := 8192
	numGpus := 4
	engine := "vllm"
	preset := "llama3-70b"

	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		if p := r.FormValue("preset"); p != "" {
			preset = p
			switch p {
			case "llama3-70b":
				modelName = "meta-llama/Llama-3-70b-Instruct"
			case "deepseek-v3":
				modelName = "deepseek-ai/DeepSeek-V3"
			case "qwen-2.5-72b":
				modelName = "Qwen/Qwen2.5-72B-Instruct"
			case "mistral-large":
				modelName = "mistralai/Mistral-Large-Instruct-2407"
			case "custom":
				if m := r.FormValue("modelName"); m != "" {
					modelName = m
				}
			}
		} else if m := r.FormValue("modelName"); m != "" {
			modelName = m
		}
		if l := r.FormValue("maxModelLen"); l != "" {
			if n, err := strconv.Atoi(l); err == nil {
				maxModelLen = n
			}
		}
		if g := r.FormValue("numGpus"); g != "" {
			if n, err := strconv.Atoi(g); err == nil {
				numGpus = n
			}
		}
		if e := r.FormValue("engine"); e != "" {
			engine = e
		}
	}

	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.90,
		NumGpus:              numGpus,
		MaxNumSeqs:           32,
		DtypeBytes:           2,
	}
	a := Analyze(cfg)
	command := buildCommand(engine, modelName, cfg, a)

	vllmSel, sglSel := "", ""
	if engine == "sglang" {
		sglSel = "selected"
	} else {
		vllmSel = "selected"
	}

	sel := map[string]string{}
	sel[preset] = "selected"

	oomText := "false"
	oomColor := "#34d399"
	if a.OOM {
		oomText = "true"
		oomColor = "#f43f5e"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Pusula Serve</title>
<style>
body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background-color: #070b14;
    color: #f8fafc;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
}
.app-container {
    width: 100%%;
    max-width: 900px;
    padding: 40px 20px;
}
.header-brand {
    font-size: 32px;
    font-weight: 800;
    color: #38bdf8;
    margin-bottom: 30px;
    letter-spacing: -0.5px;
}
.form-grid {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 35px;
}
.field-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
label {
    font-size: 14px;
    font-weight: 500;
    color: #cbd5e1;
}
select, input {
    width: 100%%;
    padding: 12px 16px;
    background-color: #0b132b;
    border: 1px solid #1e293b;
    color: #fff;
    border-radius: 8px;
    font-size: 14px;
    box-sizing: border-box;
}
select:focus, input:focus {
    outline: none;
    border-color: #38bdf8;
}
.results-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    font-size: 15px;
    font-weight: 600;
    color: #94a3b8;
}
.calculated-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #38bdf8;
    font-size: 14px;
}
.cards-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 12px;
    margin-bottom: 30px;
}
.metric-card {
    background-color: #0b132b;
    border: 1px solid #1e293b;
    border-radius: 10px;
    padding: 16px 12px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-height: 85px;
}
.card-top {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: #94a3b8;
    font-weight: 500;
}
.card-val {
    font-size: 18px;
    font-weight: 700;
    color: #38bdf8;
    margin-top: 8px;
}
.cmd-box {
    background-color: #0b132b;
    border: 1px solid #1e293b;
    border-radius: 10px;
    padding: 20px;
    margin-top: 20px;
}
.cmd-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #94a3b8;
    font-size: 14px;
    margin-bottom: 12px;
}
pre {
    margin: 0;
    background: #04060b;
    padding: 14px;
    border-radius: 6px;
    color: #38bdf8;
    overflow-x: auto;
    font-family: monospace;
    font-size: 13px;
    line-height: 1.5;
}
</style>
</head>
<body>

<div class="app-container">
    <div class="header-brand">Pusula Serve</div>

    <form method="POST">
        <div class="form-grid">
            <div class="field-group">
                <label>Engine</label>
                <select name="engine">
                    <option value="vllm" %s>vLLM</option>
                    <option value="sglang" %s>SGLang</option>
                </select>
            </div>
            <div class="field-group">
                <label>Preset</label>
                <select name="preset" onchange="this.form.submit()">
                    <option value="llama3-70b" %s>Llama-3-70B</option>
                    <option value="deepseek-v3" %s>DeepSeek-V3</option>
                    <option value="qwen-2.5-72b" %s>Qwen-2.5-72B</option>
                    <option value="mistral-large" %s>Mistral-Large</option>
                    <option value="custom" %s>Custom</option>
                </select>
            </div>
            <div class="field-group">
                <label>max-model-len</label>
                <input type="number" name="maxModelLen" value="%d" onchange="this.form.submit()">
            </div>
            <div class="field-group">
                <label>GPU count</label>
                <select name="numGpus" onchange="this.form.submit()">
                    <option value="1" %s>1</option>
                    <option value="2" %s>2</option>
                    <option value="4" %s>4</option>
                    <option value="8" %s>8</option>
                </select>
            </div>
        </div>
    </form>

    <div class="results-header">
        <span>Results</span>
        <div class="calculated-badge">
            <span>Calculated</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M20 6L9 17l-5-5"/></svg>
        </div>
    </div>

    <div class="cards-grid">
        <div class="metric-card">
            <div class="card-top">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
                <span>Weights</span>
            </div>
            <div class="card-val">%.1f GB</div>
        </div>

        <div class="metric-card">
            <div class="card-top">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 6h16M4 12h16m-7 6h7"/></svg>
                <span>KV Cache</span>
            </div>
            <div class="card-val">%.1f GB</div>
        </div>

        <div class="metric-card">
            <div class="card-top">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M9 9h6v6H9z"/></svg>
                <span>Per GPU</span>
            </div>
            <div class="card-val">%.1f GB</div>
        </div>

        <div class="metric-card">
            <div class="card-top">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4m0 4h.01"/></svg>
                <span>OOM</span>
            </div>
            <div class="card-val" style="color: %s;">%s</div>
        </div>

        <div class="metric-card">
            <div class="card-top">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><path d="M8.59 13.51l6.83 3.98m-.01-10.98l-6.83 3.98"/></svg>
                <span>TP / PP</span>
            </div>
            <div class="card-val" style="font-size: 15px; margin-top: 10px;">TP: %d<br>PP: %d</div>
        </div>
    </div>

    <div class="cmd-box">
        <div class="cmd-header">
            <span>&gt;_ Generated Command</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
        </div>
        <pre>%s</pre>
    </div>
</div>

</body>
</html>`,
		vllmSel, sglSel,
		sel["llama3-70b"], sel["deepseek-v3"], sel["qwen-2.5-72b"], sel["mistral-large"], sel["custom"],
		maxModelLen,
		map[bool]string{true: "selected", false: ""}[numGpus == 1],
		map[bool]string{true: "selected", false: ""}[numGpus == 2],
		map[bool]string{true: "selected", false: ""}[numGpus == 4],
		map[bool]string{true: "selected", false: ""}[numGpus == 8],
		a.WeightGB,
		a.KVGB,
		a.PerGPUGB,
		oomColor, oomText,
		a.TP, a.PP,
		command,
	)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	http.HandleFunc("/api/v1/optimize", handleAPI)
	fmt.Println("Pusula Serve :8080 adresinde calisiyor...")
	_ = http.ListenAndServe(":8080", nil)
}
