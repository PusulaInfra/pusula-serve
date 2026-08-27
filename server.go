package main

import (
	"fmt"
	"net/http"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	command := "vllm serve meta-llama/Llama-3-70b-Instruct --max-model-len 32768 --gpu-memory-utilization 0.90 --tensor-parallel-size 4"
	
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
				h1 {
					color: #38bdf8;
					font-size: 24px;
					margin-top: 0;
				}
				.badge {
					background: #0369a1;
					color: #e0f2fe;
					padding: 4px 10px;
					border-radius: 20px;
					font-size: 12px;
					font-weight: bold;
				}
				.card {
					background: #1f2937;
					border-left: 4px solid #38bdf8;
					padding: 15px;
					margin: 20px 0;
					border-radius: 4px;
				}
				pre {
					background: #030712;
					padding: 15px;
					border-radius: 6px;
					color: #34d399;
					overflow-x: auto;
					font-family: monospace;
				}
				p { color: #94a3b8; }
			</style>
		</head>
		<body>
			<div class="container">
				<span class="badge">vLLM & SGLang Optimizer</span>
				<h1>Pusula Serve Copilot</h1>
				<p>Yapay zeka modellerinizi cluster üzerinde en verimli şekilde sunmanız için optimize edilmiş konfigürasyon motoru.</p>
				
				<div class="card">
					<h3>Önerilen Üretim Komutu:</h3>
					<pre>%s</pre>
				</div>
			</div>
		</body>
		</html>
	`, command)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	fmt.Println("Pusula Serve gelistirilmis arayüzü ile baslatiliyor...")
	http.ListenAndServe(":8080", nil)
}
