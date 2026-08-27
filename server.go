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
			"python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --mem-fraction-static %.2f --port 30000",
			model, cfg.MaxModelLen, a.TP, cfg.GpuMemoryUtilization)
	}
	return fmt.Sprintf(
		"vllm serve %s --max-model-len %d --tensor-parallel-size %d --pipeline-parallel-size %d --gpu-memory-utilization %.2f --max-num-seqs %d",
		model, cfg.MaxModelLen, a.TP, a.PP, cfg.GpuMemoryUtilization, cfg.MaxNumSeqs)
}

func buildK8s(engine, model string, cfg ServingConfig, a Analysis) string {
	image := "vllm/vllm-openai:v0.8.5"
	containerCmd := fmt.Sprintf(`["vllm","serve","%s","--max-model-len","%d","--tensor-parallel-size","%d"]`, model, cfg.MaxModelLen, a.TP)
	if engine == "sglang" {
		image = "lmsysorg/sglang:latest"
		containerCmd = fmt.Sprintf(`["python3","-m","sglang.launch_server","--model-path","%s","--tp-size","%d"]`, model, a.TP)
	}
	return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: pusula-serve
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pusula-serve
  template:
    spec:
      containers:
      - name: llm
        image: %s
        command: %s
        resources:
          limits:
            nvidia.com/gpu: %d
        volumeMounts:
        - name: shm
          mountPath: /dev/shm
      volumes:
      - name: shm
        emptyDir:
          medium: Memory
          sizeLimit: 8Gi
`, image, containerCmd, cfg.NumGpus)
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
		MaxModelLen:          queryInt(r, "max_model_len", 32768),
		NumGpus:              queryInt(r, "num_gpus", 4),
		MaxNumSeqs:           queryInt(r, "max_num_seqs", 16),
		GpuMemoryUtilization: 0.90,
		DtypeBytes:           2,
	}
	a := Analyze(cfg)
	hourly := 2.15 * float64(cfg.NumGpus)
	cmd := buildCommand(engine, modelName, cfg, a)
	json.NewEncoder(w).Encode(APIResponse{
		ModelName: modelName, Engine: engine,
		WeightGB: a.WeightGB, KVGB: a.KVGB, VramGB: a.TotalGB, PerGPUGB: a.PerGPUGB,
		TP: a.TP, PP: a.PP, OOM: a.OOM,
		MonthlyCostUSD: hourly * 24 * 30,
		Command:        cmd,
		Recommendation: a.Recommendation,
	})
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	modelName := "meta-llama/Llama-3-70b-Instruct"
	maxModelLen := 32768
	numGpus := 4
	maxNumSeqs := 16
	engine := "vllm"
	lang := "tr"
	preset := "llama3-70b"

	if r.Method == http.MethodPost {
		r.ParseForm()
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
			if parsedLen, err := strconv.Atoi(l); err == nil {
				maxModelLen = parsedLen
			}
		}
		if g := r.FormValue("numGpus"); g != "" {
			if parsedGpu, err := strconv.Atoi(g); err == nil {
				numGpus = parsedGpu
			}
		}
		if s := r.FormValue("maxNumSeqs"); s != "" {
			if parsedSeq, err := strconv.Atoi(s); err == nil {
				maxNumSeqs = parsedSeq
			}
		}
		if e := r.FormValue("engine"); e != "" {
			engine = e
		}
		if lg := r.FormValue("lang"); lg != "" {
			lang = lg
		}
	}

	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.90,
		NumGpus:              numGpus,
		MaxNumSeqs:           maxNumSeqs,
		DtypeBytes:           2,
	}

	a := Analyze(cfg)
	vram := a.TotalGB
	recommendation := a.Recommendation
	tpSize, ppSize := a.TP, a.PP

	hourlyCostPerGpu := 2.15
	totalHourlyCost := float64(numGpus) * hourlyCostPerGpu
	monthlyCost := totalHourlyCost * 24 * 30

	command := buildCommand(engine, modelName, cfg, a)
	k8sYaml := buildK8s(engine, modelName, cfg, a)

	var tTitle, tNavDashboard, tNavModels, tNavDocs, tNavGithub, tBadge, tDesc, tLangLbl, tEngLbl, tPresetLbl, tModelLbl, tLenLbl, tGpuLbl, tBtn, tCard1Title, tVramText, tRecText, tCostTitle, tCard2Title, tCard3Title string

	if lang == "en" {
		tTitle = "Pusula Serve - Realist LLM Serving MVP"
		tNavDashboard = "Dashboard"
		tNavModels = "AI Models"
		tNavDocs = "API & Docs"
		tNavGithub = "GitHub Repo"
		tBadge = "MVP Active (No-BS Engine)"
		tDesc = "Model preset + context + GPU count -> realistic VRAM/KV estimation, TP/PP selection, and launch command."
		tLangLbl = "Language / Dil:"
		tEngLbl = "Serving Engine:"
		tPresetLbl = "Model Preset:"
		tModelLbl = "Custom Model Name (HuggingFace):"
		tLenLbl = "Context Window (Tokens):"
		tGpuLbl = "Total GPU Count:"
		tBtn = "Run Realistic Analysis"
		tCard1Title = "Hardware & Memory Breakdown:"
		tVramText = "Memory Allocation Analysis:"
		tRecText = "Optimal Cluster Topology:"
		tCostTitle = "Cloud Infrastructure Cost Estimation:"
		tCard2Title = "Generated Serving Command"
		tCard3Title = "Kubernetes (K8s) Deployment Manifest"
	} else {
		tTitle = "Pusula Serve - Gerçekçi LLM Serving MVP"
		tNavDashboard = "Panel (Dashboard)"
		tNavModels = "Modeller"
		tNavDocs = "Dokümantasyon"
		tNavGithub = "GitHub Kaynağı"
		tBadge = "MVP Aktif (Gerçekçi Hesaplama Motoru)"
		tDesc = "Model preset + context + GPU sayısı ile gerçekçi VRAM/KV analizi, TP/PP seçimi ve başlatma komutu."
		tLangLbl = "Dil / Language:"
		tEngLbl = "Serving Motoru:"
		tPresetLbl = "Hazır Model Seçimi (Preset):"
		tModelLbl = "Özel Model Adı (HuggingFace):"
		tLenLbl = "Context Window (Token):"
		tGpuLbl = "Toplam GPU Sayısı:"
		tBtn = "Gerçekçi Analizi Çalıştır"
		tCard1Title = "Donanım & Bellek Dökümü:"
		tVramText = "Bellek Analizi:"
		tRecText = "Önerilen Cluster Topolojisi:"
		tCostTitle = "Bulut Altyapı Maliyet Tahmini:"
		tCard2Title = "Üretilen Başlatma Komutu"
		tCard3Title = "Kubernetes (K8s) Deployment Manifestosu"
	}

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="%s">
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
			<style>
				body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background-color: #090d16; color: #e2e8f0; margin: 0; padding: 0; }
				.navbar { background: #111827; border-bottom: 1px solid #1f2937; padding: 15px 40px; display: flex; justify-content: space-between; align-items: center; }
				.navbar .logo { color: #38bdf8; font-weight: bold; font-size: 18px; text-decoration: none; }
				.navbar .nav-links { display: flex; gap: 20px; }
				.navbar .nav-links a { color: #94a3b8; text-decoration: none; font-size: 14px; font-weight: 500; }
				.navbar .nav-links a:hover { color: #38bdf8; }
				.main-wrapper { display: flex; justify-content: center; padding: 40px; }
				.container { width: 100%%; max-width: 850px; background: #111827; border: 1px solid #1f2937; padding: 30px; border-radius: 12px; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3); }
				h1 { color: #38bdf8; font-size: 24px; margin-top: 0; }
				.badge { background: #0369a1; color: #e0f2fe; padding: 4px 10px; border-radius: 20px; font-size: 12px; font-weight: bold; }
				.card { background: #1f2937; border-left: 4px solid #38bdf8; padding: 15px; margin: 20px 0; border-radius: 4px; }
				.metric { font-size: 16px; color: #34d399; font-weight: bold; }
				.cost { font-size: 18px; color: #f59e0b; font-weight: bold; }
				pre { background: #030712; padding: 15px; border-radius: 6px; color: #38bdf8; overflow-x: auto; font-family: monospace; font-size: 13px; }
				p { color: #94a3b8; }
				label { display: block; margin-top: 12px; color: #cbd5e1; font-weight: 500; font-size: 14px; }
				input, select { width: 100%%; padding: 10px; margin-top: 5px; background: #030712; border: 1px solid #374151; color: #fff; border-radius: 6px; box-sizing: border-box; }
				button { margin-top: 20px; background: #0284c7; color: white; border: none; padding: 12px 20px; font-size: 16px; font-weight: bold; border-radius: 6px; cursor: pointer; width: 100%%; }
				button:hover { background: #0369a1; }
			</style>
		</head>
		<body>
			<nav class="navbar">
				<a href="#" class="logo">🧭 Pusula Serve MVP</a>
				<div class="nav-links">
					<a href="#">%s</a>
					<a href="#">%s</a>
					<a href="#">%s</a>
					<a href="https://github.com/ccyhun67-gif/pusula-serve" target="_blank">%s</a>
				</div>
			</nav>
			<div class="main-wrapper">
				<div class="container">
					<span class="badge">%s</span>
					<h1>Pusula Serve Core</h1>
					<p>%s</p>
					<form method="POST">
						<label>%s</label>
						<select name="lang" onchange="this.form.submit()">
							<option value="tr" %s>Türkçe (TR)</option>
							<option value="en" %s>English (EN)</option>
						</select>
						<label>%s</label>
						<select name="engine">
							<option value="vllm" %s>vLLM (High-Throughput)</option>
							<option value="sglang" %s>SGLang (RadixAttention & Low-Latency)</option>
						</select>
						<label>%s</label>
						<select name="preset" onchange="this.form.submit()">
							<option value="llama3-70b" %s>Llama-3-70B-Instruct</option>
							<option value="deepseek-v3" %s>DeepSeek-V3 (MoE)</option>
							<option value="qwen-2.5-72b" %s>Qwen-2.5-72B-Instruct</option>
							<option value="mistral-large" %s>Mistral-Large-Instruct</option>
							<option value="custom" %s>Özel Model Gir (Custom Model)</option>
						</select>
						<label>%s</label>
						<input type="text" name="modelName" value="%s">
						<label>%s</label>
						<input type="number" name="maxModelLen" value="%d">
						<label>%s</label>
						<input type="number" name="numGpus" value="%d">
						<button type="submit">%s</button>
					</form>
					<div class="card">
						<h3>%s</h3>
						<p class="metric">Weights: %.1f GB · KV: %.1f GB · Per GPU: %.1f GB (Usable: %.1f GB)</p>
						<p>%s (TP: %d, PP: %d) — OOM Status: <b>%v</b></p>
						<p style="color: #38bdf8; margin-top: 8px;">💡 <b>Analiz:</b> %s</p>
					</div>
					<div class="card" style="border-left-color: #f59e0b;">
						<h3 style="color: #fbbf24;">%s</h3>
						<p class="cost">Saatlik Tahmini Maliyet: $%.2f / saat</p>
						<p class="cost">Aylık Tahmini Bulut Maliyeti (7/24): $%.2f / ay (%d x GPU)</p>
					</div>
					<div class="card">
						<h3>%s (%s):</h3>
						<pre>%s</pre>
					</div>
					<div class="card" style="border-left-color: #10b981;">
						<h3 style="color: #34d399;">%s</h3>
						<pre>%s</pre>
					</div>
				</div>
			</div>
		</body>
		</html>
	`, 
	lang, tTitle, 
	tNavDashboard, tNavModels, tNavDocs, tNavGithub,
	tBadge, tDesc, 
	tLangLbl, 
	func() string { if lang == "tr" { return "selected" } ; return "" }(),
	func() string { if lang == "en" { return "selected" } ; return "" }(),
	tEngLbl,
	func() string { if engine == "vllm" { return "selected" } ; return "" }(),
	func() string { if engine == "sglang" { return "selected" } ; return "" }(),
	tPresetLbl,
	func() string { if preset == "llama3-70b" { return "selected" } ; return "" }(),
	func() string { if preset == "deepseek-v3" { return "selected" } ; return "" }(),
	func() string { if preset == "qwen-2.5-72b" { return "selected" } ; return "" }(),
	func() string { if preset == "mistral-large" { return "selected" } ; return "" }(),
	func() string { if preset == "custom" { return "selected" } ; return "" }(),
	tModelLbl, modelName,
	tLenLbl, maxModelLen,
	tGpuLbl, numGpus,
	tBtn,
	tCard1Title, a.WeightGB, a.KVGB, a.PerGPUGB, a.UsablePerGPU, tRecText, tpSize, ppSize, a.OOM, recommendation,
	tCostTitle, totalHourlyCost, monthlyCost, numGpus,
	tCard2Title, engine, command,
	tCard3Title, k8sYaml)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	http.HandleFunc("/api/v1/optimize", handleAPI)
	fmt.Println("Pusula Serve Gerçekçi MVP Sunucusu başlatılıyor...")
	http.ListenAndServe(":8080", nil)
}
