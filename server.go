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

	// Kullanıcı formdan yeni değerler gönderdiyse onları alıyoruz
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
	}

	// Hesaplama motorunu çalıştırıyoruz
	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.90,
		NumGpus:              numGpus,
	}

	vram, recommendation := CalculateVRAMAndCost(cfg)
	command := fmt.Sprintf("vllm serve %s --max-model-len %d --tensor-parallel-size %d --gpu-memory-utilization 0.90", modelName, maxModelLen, numGpus)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="tr">
		<head>
			<meta charset="UTF-8">
			<title>Pusula Serve - LLM Serving Copilot</title>
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
				<span class="badge">İnteraktif LLM Optimizasyon Paneli</span>
				<h1>Pusula Serve Copilot</h1>
				<p>Model ayarlarını yapılandırın, VRAM ihtiyacını ve hazır vLLM komutunu anında alın.</p>
				
				<form method="POST">
					<label>Model Adı (HuggingFace Repo):</label>
					<input type="text" name="modelName" value="%s">

					<label>Max Model Uzunluğu (Context Window - Token):</label>
					<input type="number" name="maxModelLen" value="%d">

					<label>GPU Sayısı (Tensor Parallelism):</label>
					<input type="number" name="numGpus" value="%d">

					<button type="submit">Optimize Et ve Komut Üret</button>
				</form>

				<div class="card">
					<h3>Donanım & Bellek Analizi:</h3>
					<p class="metric">Tahmini VRAM Tüketimi: %.2f GB</p>
					<p>Cluster Tavsiyesi: %s</p>
				</div>

				<div class="card">
					<h3>Üretilen vLLM Komutu:</h3>
					<pre>%s</pre>
				</div>
			</div>
		</body>
		</html>
	`, modelName, maxModelLen, numGpus, vram, recommendation, command)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	fmt.Println("Pusula Serve interaktif web sunucusu baslatiliyor...")
	http.ListenAndServe(":8080", nil)
}
