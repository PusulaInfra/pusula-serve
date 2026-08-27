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
	lang := "tr" // Varsayılan dil Türkçe

	// Kullanıcı formdan yeni değerler gönderdiyse alıyoruz
	if r.Method == http.MethodPost {
		r.ParseForm()
		if m := r.FormValue("modelName"); m != "" {
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

	// Hesaplama motorunu çalıştırıyoruz
	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.90,
		NumGpus:              numGpus,
	}

	vram, recommendation := CalculateVRAMAndCost(cfg)
	
	// Seçilen motora göre komut üretimi
	var command string
	if engine == "sglang" {
		command = fmt.Sprintf("python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --port 30000", modelName, maxModelLen, numGpus)
	} else {
		command = fmt.Sprintf("vllm serve %s --max-model-len %d --tensor-parallel-size %d --gpu-memory-utilization 0.90", modelName, maxModelLen, numGpus)
	}

	// Metinlerin sözlüğü (TR / EN)
	var tTitle, tBadge, tDesc, tLangLbl, tEngLbl, tModelLbl, tLenLbl, tGpuLbl, tBtn, tCard1Title, tVramText, tRecText, tCard2Title string

	if lang == "en" {
		tTitle = "Pusula Serve - LLM Serving Copilot"
		tBadge = "Multi-Engine & Multi-Language Active"
		tDesc = "Configure your serving engine, parameters, and get optimized commands instantly."
		tLangLbl = "Language / Dil:"
		tEngLbl = "Serving Engine:"
		tModelLbl = "Model Name (HuggingFace Repo):"
		tLenLbl = "Max Model Length (Context Window - Tokens):"
		tGpuLbl = "GPU Count (Tensor Parallelism):"
		tBtn = "Optimize & Generate Command"
		tCard1Title = "Hardware & Memory Analysis:"
		tVramText = "Estimated VRAM Usage:"
		tRecText = "Cluster Recommendation:"
		tCard2Title = "Generated Startup Command"
	} else {
		tTitle = "Pusula Serve - LLM Serving Copilot"
		tBadge = "Multi-Engine & Çok Dilli Destek Aktif"
		tDesc = "Serving motorunu seçin, parametreleri yapılandırın ve optimize komutu anında alın."
		tLangLbl = "Dil / Language:"
		tEngLbl = "Serving Motoru:"
		tModelLbl = "Model Adı (HuggingFace Repo):"
		tLenLbl = "Max Model Uzunluğu (Context Window - Token):"
		tGpuLbl = "GPU Sayısı (Tensor Parallelism):"
		tBtn = "Optimize Et ve Komut Üret"
		tCard1Title = "Donanım & Bellek Analizi:"
		tVramText = "Tahmini VRAM Tüketimi:"
		tRecText = "Cluster Tavsiyesi:"
		tCard2Title = "Üretilen Başlatma Komutu"
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
					max-width: 800px;
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
				pre { background: #030712; padding: 15px; border-radius: 6px; color: #38bdf8; overflow-x: auto; font-family: monospace; }
				p { color: #94a3b8; }
				label { display: block; margin-top: 10px; color: #cbd5e1; font-weight: 500; }
				input, select { width: 100%%; padding: 10px; margin-top: 5px; background: #030712; border: 1px solid #374151; color: #fff; border-radius: 6px; box-sizing: border-box; }
				button { margin-top: 20px; background: #0284c7; color: white; border: none; padding: 12px 20px; font-size: 16px; font-weight: bold; border-radius: 6px; cursor: pointer; width: 100%%; transition: background 0.2s; }
				button:hover { background: #0369a1; }
			</style>
		</head>
		<body>
			<div class="container">
				<span class="badge">%s</span>
				<h1>Pusula Serve Copilot</h1>
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
					<p>%s %s</p>
				</div>

				<div class="card">
					<h3>%s (%s):</h3>
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
	tModelLbl, modelName,
	tLenLbl, maxModelLen,
	tGpuLbl, numGpus,
	tBtn,
	tCard1Title, tVramText, vram, tRecText, recommendation,
	tCard2Title, engine, command)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	fmt.Println("Pusula Serve Çok Dilli web sunucusu baslatiliyor...")
	http.ListenAndServe(":8080", nil)
}
