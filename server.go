package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	// Varsayılan değerler
	modelName := "meta-llama/Llama-3-70b-Instruct"
	maxModelLen := 32768
	numGpus := 4
	engine := "vllm"
	lang := "tr"
	preset := "llama3-70b"

	// Formdan gelen değerleri alıyoruz
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
	}

	vram, recommendation := CalculateVRAMAndCost(cfg)
	
	// Akıllı TP ve PP Hesaplama Mantığı
	tpSize := numGpus
	ppSize := 1
	if numGpus >= 8 && maxModelLen > 65536 {
		tpSize = 4
		ppSize = 2
	}

	// Bulut Maliyet Hesaplama Simülasyonu (GPU başına ortalama saatlik $2.15 baz alınıyor)
	hourlyCostPerGpu := 2.15
	totalHourlyCost := float64(numGpus) * hourlyCostPerGpu
	monthlyCost := totalHourlyCost * 24 * 30

	// Komut üretimi
	var command, k8sYaml string
	if engine == "sglang" {
		command = fmt.Sprintf("python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --port 30000", modelName, maxModelLen, tpSize)
	} else {
		command = fmt.Sprintf("vllm serve %s --max-model-len %d --tensor-parallel-size %d --pipeline-parallel-size %d --gpu-memory-utilization 0.90", modelName, maxModelLen, tpSize, ppSize)
	}

	// Kubernetes YAML Manifest Örneği
	k8sYaml = fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: llm-serving-%s
  labels:
    app: pusula-serve
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pusula-serve
  template:
    metadata:
      labels:
        app: pusula-serve
    spec:
      containers:
      - name: llm-engine
        image: vllm/vllm-openai:latest
        command: ["vllm", "serve", "%s", "--max-model-len", "%d", "--tensor-parallel-size", "%d"]
        resources:
          limits:
            nvidia.com/gpu: "%d"`, engine, modelName, maxModelLen, tpSize, numGpus)

	// Çok dilli sözlük
	var tTitle, tBadge, tDesc, tLangLbl, tEngLbl, tPresetLbl, tModelLbl, tLenLbl, tGpuLbl, tBtn, tCard1Title, tVramText, tRecText, tCostTitle, tCard2Title, tCard3Title string

	if lang == "en" {
		tTitle = "Pusula Serve - Enterprise LLM SaaS Platform"
		tBadge = "Elite Edition Active (Cost Calculator & K8s)"
		tDesc = "Simulate multi-GPU cluster economics, optimize deployment topology, and export production-ready manifests."
		tLangLbl = "Language / Dil:"
		tEngLbl = "Serving Engine:"
		tPresetLbl = "Model Preset:"
		tModelLbl = "Custom Model Name (HuggingFace):"
		tLenLbl = "Context Window (Tokens):"
		tGpuLbl = "Total GPU Count:"
		tBtn = "Run Full Enterprise Analysis"
		tCard1Title = "Hardware & Cluster Intelligence:"
		tVramText = "Estimated VRAM Consumption:"
		tRecText = "Optimal Cluster Topology:"
		tCostTitle = "Cloud Infrastructure Cost Estimation:"
		tCard2Title = "Generated Serving Command"
		tCard3Title = "Kubernetes (K8s) Deployment Manifest"
	} else {
		tTitle = "Pusula Serve - Kurumsal LLM SaaS Platformu"
		tBadge = "Elite Sürüm Aktif (Maliyet Hesaplayıcı & K8s)"
		tDesc = "Çoklu GPU cluster ekonomisini simüle edin, dağılım topolojisini optimize edin ve üretime hazır manifestolar üretin."
		tLangLbl = "Dil / Language:"
		tEngLbl = "Serving Motoru:"
		tPresetLbl = "Hazır Model Seçimi (Preset):"
		tModelLbl = "Özel Model Adı (HuggingFace):"
		tLenLbl = "Context Window (Token):"
		tGpuLbl = "Toplam GPU Sayısı:"
		tBtn = "Tam Kurumsal Analizi Çalıştır"
		tCard1Title = "Donanım & Cluster Zekası:"
		tVramText = "Tahmini VRAM Tüketimi:"
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
				body {
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
					background-color: #090d16;
					color: #e2e8f0;
					margin: 0;
					padding: 40px;
					display: flex;
					justify-content: center;
				}
				.container {
					width: 100%%;
					max-width: 850px;
					background: #111827;
					border: 1px solid #1f2937;
					padding: 30px;
					border-radius: 12px;
					box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3);
				}
				h1 { color: #38bdf8; font-size: 24px; margin-top: 0; }
				.badge { background: #0369a1; color: #e0f2fe; padding: 4px 10px; border-radius: 20px; font-size: 12px; font-weight: bold; }
				.card { background: #1f2937; border-left: 4px solid #38bdf8; padding: 15px; margin: 20px 0; border-radius: 4px; }
				.metric { font-size: 18px; color: #34d399; font-weight: bold; }
				.cost { font-size: 18px; color: #f59e0b; font-weight: bold; }
				pre { background: #030712; padding: 15px; border-radius: 6px; color: #38bdf8; overflow-x: auto; font-family: monospace; font-size: 13px; }
				p { color: #94a3b8; }
				label { display: block; margin-top: 12px; color: #cbd5e1; font-weight: 500; font-size: 14px; }
				input, select { width: 100%%; padding: 10px; margin-top: 5px; background: #030712; border: 1px solid #374151; color: #fff; border-radius: 6px; box-sizing: border-box; }
				button { margin-top: 20px; background: #0284c7; color: white; border: none; padding: 12px 20px; font-size: 16px; font-weight: bold; border-radius: 6px; cursor: pointer; width: 100%%; transition: background 0.2s; }
				button:hover { background: #0369a1; }
			</style>
		</head>
		<body>
			<div class="container">
				<span class="badge">%s</span>
				<h1>Pusula Serve Elite SaaS</h1>
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
					<p class="metric">%s %.2f GB</p>
					<p>%s (TP: %d, PP: %d) — %s</p>
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
		</body>
		</html>
	`, 
	lang, tTitle, tBadge, tDesc, 
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
	tCard1Title, tVramText, vram, tRecText, tpSize, ppSize, recommendation,
	tCostTitle, totalHourlyCost, monthlyCost, numGpus,
	tCard2Title, engine, command,
	tCard3Title, k8sYaml)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	fmt.Println("Pusula Serve Elite SaaS Sürümü başlatılıyor...")
	http.ListenAndServe(":8080", nil)
}
