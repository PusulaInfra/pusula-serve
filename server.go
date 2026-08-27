package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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

func getCloudRate(provider string, gpu string, gpus int) float64 {
	rates := map[string]map[string]float64{
		"lambda": {"H100": 2.49, "A100": 1.29, "L40S": 0.79},
		"runpod": {"H100": 2.29, "A100": 1.19, "L40S": 0.74},
		"aws":    {"H100": 4.10, "A100": 3.67, "L40S": 1.80},
		"gcp":    {"H100": 3.99, "A100": 3.50, "L40S": 1.75},
	}
	if p, ok := rates[provider]; ok {
		if rate, exists := p[gpu]; exists {
			return rate * float64(gpus)
		}
	}
	return 2.49 * float64(gpus)
}

func buildProdCommand(engine, model string, cfg ServingConfig, a Analysis) string {
	if engine == "sglang" {
		return fmt.Sprintf(
			"python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --mem-fraction-static %.2f --port 8000",
			model, cfg.MaxModelLen, a.TP, cfg.GpuMemoryUtilization)
	}
	return fmt.Sprintf(
		"vllm serve %s --max-model-len %d --tensor-parallel-size %d --pipeline-parallel-size %d --gpu-memory-utilization %.2f --max-num-seqs %d --enable-prefix-caching",
		model, cfg.MaxModelLen, a.TP, a.PP, cfg.GpuMemoryUtilization, cfg.MaxNumSeqs)
}

func handleEnterpriseConfig(w http.ResponseWriter, r *http.Request) {
	modelName := "meta-llama/Llama-3-70b-Instruct"
	maxModelLen := 16384
	numGpus := 4
	engine := "vllm"
	provider := "lambda"
	gpuType := "H100"
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
			case "llama-405b":
				modelName = "meta-llama/Llama-3.1-405B-Instruct"
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
		if pr := r.FormValue("provider"); pr != "" {
			provider = pr
		}
		if gt := r.FormValue("gpuType"); gt != "" {
			gpuType = gt
		}
	}

	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.92,
		NumGpus:              numGpus,
		MaxNumSeqs:           32,
		DtypeBytes:           2,
	}
	a := Analyze(cfg)
	hourlyRate := getCloudRate(provider, gpuType, numGpus)
	monthlyCost := hourlyRate * 24 * 30
	command := buildProdCommand(engine, modelName, cfg, a)

	selEngine := map[string]string{engine: "selected"}
	selPreset := map[string]string{preset: "selected"}
	selProvider := map[string]string{provider: "selected"}
	selGpuType := map[string]string{gpuType: "selected"}

	oomStatus := "OPTIMIZED"
	statusColor := "#34d399"
	if a.OOM {
		oomStatus = "OOM CRITICAL"
		statusColor = "#f43f5e"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Pusula Serve // AI Infrastructure & Topology Intelligence</title>
<style>
:root {
    --bg-main: #02040a;
    --surface: #090e1a;
    --border: #1e293b;
    --accent: #38bdf8;
    --accent-glow: rgba(56, 189, 248, 0.15);
    --text-primary: #f8fafc;
    --text-muted: #94a3b8;
}
body {
    font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background-color: var(--bg-main);
    background-image: 
        radial-gradient(circle at 500px -200px, rgba(56, 189, 248, 0.15), transparent 700px),
        linear-gradient(to right, rgba(30, 41, 59, 0.2) 1px, transparent 1px),
        linear-gradient(to bottom, rgba(30, 41, 59, 0.2) 1px, transparent 1px);
    background-size: 100%% 100%%, 40px 40px, 40px 40px;
    color: var(--text-primary);
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
    min-height: 100vh;
}
.container {
    width: 100%%;
    max-width: 1100px;
    padding: 50px 24px;
    box-sizing: border-box;
}
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border);
    padding-bottom: 24px;
    margin-bottom: 35px;
    background: rgba(9, 14, 26, 0.7);
    backdrop-filter: blur(10px);
}
.logo-area {
    display: flex;
    align-items: center;
    gap: 12px;
}
.logo-title {
    font-size: 22px;
    font-weight: 800;
    color: var(--accent);
    letter-spacing: -0.5px;
}
.version-badge {
    background: var(--accent-glow);
    color: var(--accent);
    border: 1px solid rgba(56, 189, 248, 0.3);
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
}
.dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    margin-bottom: 30px;
}
.card {
    background-color: var(--surface);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 24px;
    box-sizing: border-box;
    box-shadow: 0 20px 40px -15px rgba(0,0,0,0.7);
    backdrop-filter: blur(10px);
}
.card-title {
    font-size: 14px;
    font-weight: 700;
    text-transform: uppercase;
    color: var(--text-muted);
    letter-spacing: 0.5px;
    margin-bottom: 16px;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    margin-bottom: 16px;
}
.form-row.single {
    grid-template-columns: 1fr;
}
label {
    display: block;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 6px;
    text-transform: uppercase;
}
select, input {
    width: 100%%;
    padding: 12px 14px;
    background-color: #030712;
    border: 1px solid var(--border);
    color: #fff;
    border-radius: 8px;
    font-size: 14px;
    box-sizing: border-box;
    transition: all 0.2s;
}
select:focus, input:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 3px var(--accent-glow);
}
.metrics-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 30px;
}
.metric-card {
    background-color: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    box-shadow: 0 10px 25px -10px rgba(0,0,0,0.5);
}
.metric-label {
    font-size: 11px;
    font-weight: 700;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
}
.metric-val {
    font-size: 22px;
    font-weight: 800;
    color: var(--accent);
    margin-top: 12px;
    letter-spacing: -0.5px;
}
.terminal-card {
    background-color: var(--surface);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 24px;
    box-shadow: 0 20px 40px -15px rgba(0,0,0,0.7);
}
pre {
    margin: 0;
    background: #030712;
    padding: 16px;
    border-radius: 8px;
    color: var(--accent);
    overflow-x: auto;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 13px;
    line-height: 1.6;
    border: 1px solid var(--border);
}
</style>
</head>
<body>

<div class="container">
    <div class="header">
        <div class="logo-area">
            <span class="logo-title">🧭 Pusula Serve // Enterprise AI</span>
            <span class="version-badge">Topology Engine v3.0</span>
        </div>
        <a href="https://github.com/ccyhun67-gif/pusula-serve" target="_blank" style="color: var(--text-muted); text-decoration: none; font-size: 13px; font-weight: 600;">GitHub Repository</a>
    </div>

    <form method="POST">
        <div class="dashboard-grid">
            <div class="card">
                <div class="card-title">Infrastructure Topology</div>
                <div class="form-row">
                    <div>
                        <label>Cloud Provider</label>
                        <select name="provider" onchange="this.form.submit()">
                            <option value="lambda" %s>Lambda Labs</option>
                            <option value="runpod" %s>RunPod</option>
                            <option value="aws" %s>AWS UltraCluster</option>
                            <option value="gcp" %s>GCP AI Hypercomputer</option>
                        </select>
                    </div>
                    <div>
                        <label>GPU Accelerator</label>
                        <select name="gpuType" onchange="this.form.submit()">
                            <option value="H100" %s>NVIDIA H100 80GB</option>
                            <option value="A100" %s>NVIDIA A100 80GB</option>
                            <option value="L40S" %s>NVIDIA L40S 48GB</option>
                        </select>
                    </div>
                </div>
                <div class="form-row">
                    <div>
                        <label>Inference Engine</label>
                        <select name="engine" onchange="this.form.submit()">
                            <option value="vllm" %s>vLLM (PagedAttention)</option>
                            <option value="sglang" %s>SGLang (RadixAttention)</option>
                        </select>
                    </div>
                    <div>
                        <label>Cluster Size (GPU Count)</label>
                        <select name="numGpus" onchange="this.form.submit()">
                            <option value="1" %s>1 x GPU</option>
                            <option value="2" %s>2 x GPUs</option>
                            <option value="4" %s>4 x GPUs</option>
                            <option value="8" %s>8 x GPUs</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-title">Model Architecture & Context</div>
                <div class="form-row single">
                    <div>
                        <label>Foundation Model Preset</label>
                        <select name="preset" onchange="this.form.submit()">
                            <option value="llama3-70b" %s>Meta Llama-3-70B-Instruct</option>
                            <option value="deepseek-v3" %s>DeepSeek-V3 (MoE)</option>
                            <option value="qwen-2.5-72b" %s>Qwen-2.5-72B-Instruct</option>
                            <option value="llama-405b" %s>Meta Llama-3.1-405B-Instruct</option>
                            <option value="custom" %s>Custom HuggingFace ID</option>
                        </select>
                    </div>
                </div>
                <div class="form-row">
                    <div>
                        <label>Custom Model Path</label>
                        <input name="modelName" value="%s">
                    </div>
                    <div>
                        <label>Max Context Length (Tokens)</label>
                        <input type="number" name="maxModelLen" value="%d" onchange="this.form.submit()">
                    </div>
                </div>
            </div>
        </div>
    </form>

    <div class="metrics-row">
        <div class="metric-card">
            <div class="metric-label">Memory Footprint</div>
            <div class="metric-val">%.1f GB <span style="font-size: 12px; color: var(--text-muted);">(Wt: %.1f | KV: %.1f)</span></div>
        </div>
        <div class="metric-card">
            <div class="metric-label">VRAM Allocation / GPU</div>
            <div class="metric-val">%.1f GB</div>
        </div>
        <div class="metric-card">
            <div class="metric-label">Estimated Economics</div>
            <div class="metric-val" style="color: #fbbf24;">$%.2f/mo</div>
        </div>
        <div class="metric-card">
            <div class="metric-label">Topology & Health</div>
            <div class="metric-val" style="font-size: 16px; color: %s;">TP: %d · PP: %d<br>%s</div>
        </div>
    </div>

    <div class="terminal-card">
        <div class="card-title">&gt;_ Production Serving Launch Command</div>
        <pre>%s</pre>
    </div>
</div>

</body>
</html>`,
		selProvider["lambda"], selProvider["runpod"], selProvider["aws"], selProvider["gcp"],
		selGpuType["H100"], selGpuType["A100"], selGpuType["L40S"],
		selEngine["vllm"], selEngine["sglang"],
		map[bool]string{true: "selected", false: ""}[numGpus == 1],
		map[bool]string{true: "selected", false: ""}[numGpus == 2],
		map[bool]string{true: "selected", false: ""}[numGpus == 4],
		map[bool]string{true: "selected", false: ""}[numGpus == 8],
		selPreset["llama3-70b"], selPreset["deepseek-v3"], selPreset["qwen-2.5-72b"], selPreset["llama-405b"], selPreset["custom"],
		modelName, maxModelLen,
		a.TotalGB, a.WeightGB, a.KVGB,
		a.PerGPUGB,
		monthlyCost,
		statusColor, a.TP, a.PP, oomStatus,
		command,
	)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleEnterpriseConfig)
	fmt.Println("Pusula Serve Enterprise v3.0 running on http://localhost:8080 ...")
	_ = http.ListenAndServe(":8080", nil)
}
