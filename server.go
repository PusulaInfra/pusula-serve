package main

import (
	"fmt"
	"net/http"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	command := "vllm serve meta-llama/Llama-3-70b-Instruct --max-model-len 32768 --gpu-memory-utilization 0.90"
	
	html := fmt.Sprintf(`
		<html>
			<head><title>Pusula Serve Copilot</title></head>
			<body style="font-family: Arial; padding: 40px; background: #0f172a; color: #f8fafc;">
				<h1>Pusula Serve - LLM Serving Copilot</h1>
				<p>Optimizasyon Motoru Aktif ve Calisiyor.</p>
				<h3>Onerilen vLLM Komutu:</h3>
				<pre style="background: #1e293b; padding: 15px; border-radius: 5px; color: #38bdf8;">%s</pre>
			</body>
		</html>
	`, command)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	fmt.Println("Pusula Serve web sunucusu baslatiliyor...")
	http.ListenAndServe(":8080", nil)
}
