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
		image = "lmsysorg/sglang:v0.4.5"
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
	maxModelLen := 32768
	numGpus := 4
	engine := "vllm"
	lang := "tr"
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
		if lg := r.FormValue("lang"); lg != "" {
			lang = lg
		}
	}

	cfg := ServingConfig{
		ModelName:            modelName,
		MaxModelLen:          maxModelLen,
		GpuMemoryUtilization: 0.90,
		NumGpus:              numGpus,
		MaxNumSeqs:           16,
		DtypeBytes:           2,
	}
	a := Analyze(cfg)
	command := buildCommand(engine, modelName, cfg, a)
	k8sYaml := buildK8s(engine, modelName, cfg, a)
	hourly := 2.15 * float64(numGpus)
	monthly := hourly * 24 * 30

	trSel, enSel := "", ""
	if lang == "en" {
		enSel = "selected"
	} else {
		trSel = "selected"
	}
	vllmSel, sglSel := "", ""
	if engine == "sglang" {
		sglSel = "selected"
	} else {
		vllmSel = "selected"
	}
	sel := map[string]string{}
	sel[preset] = "selected"

	title := "Pusula Serve"
	badge := "MVP"
	oom := "OK"
	if a.OOM {
		oom = "OOM RISK"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="%s">
<head>
<meta charset="UTF-8">
<title>%s</title>
<style>
body{font-family:-apple-system,sans-serif;background:#090d16;color:#e2e8f0;margin:0}
.navbar{background:#111827;padding:15px 40px;display:flex;justify-content:space-between}
.logo{color:#38bdf8;font-weight:bold;text-decoration:none}
.wrap{display:flex;justify-content:center;padding:40px}
.container{width:100%%;max-width:850px;background:#111827;border:1px solid #1f2937;padding:30px;border-radius:12px}
.badge{background:#0369a1;padding:4px 10px;border-radius:20px;font-size:12px}
.card{background:#1f2937;border-left:4px solid #38bdf8;padding:15px;margin:20px 0}
.metric{font-size:16px;color:#34d399;font-weight:bold}
.cost{color:#f59e0b;font-weight:bold}
pre{background:#030712;padding:15px;border-radius:6px;color:#38bdf8;overflow-x:auto;font-size:13px}
label{display:block;margin-top:12px;color:#cbd5e1}
input,select{width:100%%;padding:10px;margin-top:5px;background:#030712;border:1px solid #374151;color:#fff;border-radius:6px;box-sizing:border-box}
button{margin-top:20px;background:#0284c7;color:#fff;border:none;padding:12px;font-weight:bold;border-radius:6px;width:100%%;cursor:pointer}
</style>
</head>
<body>
<nav class="navbar">
<a class="logo" href="/">Pusula Serve</a>
<a href="https://github.com/ccyhun67-gif/pusula-serve" style="color:#94a3b8">GitHub</a>
</nav>
<div class="wrap"><div class="container">
<span class="badge">%s</span>
<h1>%s</h1>
<form method="POST">
<label>Dil</label>
<select name="lang" onchange="this.form.submit()">
<option value="tr" %s>Türkçe</option>
<option value="en" %s>English</option>
</select>
<label>Engine</label>
<select name="engine">
<option value="vllm" %s>vLLM</option>
<option value="sglang" %s>SGLang</option>
</select>
<label>Preset</label>
<select name="preset">
<option value="llama3-70b" %s>Llama-3-70B</option>
<option value="deepseek-v3" %s>DeepSeek-V3</option>
<option value="qwen-2.5-72b" %s>Qwen-2.5-72B</option>
<option value="mistral-large" %s>Mistral-Large</option>
<option value="custom" %s>Custom</option>
</select>
<label>Model</label>
<input name="modelName" value="%s">
<label>max-model-len</label>
<input type="number" name="maxModelLen" value="%d">
<label>GPU count</label>
<input type="number" name="numGpus" value="%d">
<button type="submit">Analyze</button>
</form>
<div class="card">
<p class="metric">Weights %.1f GB · KV %.1f GB · Per GPU %.1f GB · %s</p>
<p>TP %d · PP %d · %s</p>
</div>
<div class="card">
<p class="cost">$%.2f/hr · $%.2f/mo (%d x GPU @ $2.15 rough)</p>
</div>
<div class="card"><h3>Command</h3><pre>%s</pre></div>
<div class="card"><h3>K8s</h3><pre>%s</pre></div>
</div></div>
</body></html>`,
		lang, title, badge, title,
		trSel, enSel, vllmSel, sglSel,
		sel["llama3-70b"], sel["deepseek-v3"], sel["qwen-2.5-72b"], sel["mistral-large"], sel["custom"],
		modelName, maxModelLen, numGpus,
		a.WeightGB, a.KVGB, a.PerGPUGB, oom,
		a.TP, a.PP, a.Recommendation,
		hourly, monthly, numGpus,
		command, k8sYaml)

	fmt.Fprint(w, html)
}

func StartServer() {
	http.HandleFunc("/", handleConfig)
	http.HandleFunc("/api/v1/optimize", handleAPI)
	fmt.Println("Pusula Serve :8080")
	_ = http.ListenAndServe(":8080", nil)
}