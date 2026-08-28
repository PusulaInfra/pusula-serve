package serve

import "fmt"

func BuildVLLM(model string, cfg ServingConfig, a Analysis) string {
	return fmt.Sprintf(
		"vllm serve %s --max-model-len %d --tensor-parallel-size %d --pipeline-parallel-size %d --gpu-memory-utilization %.2f --max-num-seqs %d --enable-prefix-caching",
		model, cfg.MaxModelLen, a.TP, a.PP, cfg.GpuMemoryUtilization, cfg.MaxNumSeqs,
	)
}

func BuildSGLang(model string, cfg ServingConfig, a Analysis) string {
	return fmt.Sprintf(
		"python3 -m sglang.launch_server --model-path %s --context-length %d --tp-size %d --mem-fraction-static %.2f --port 8000",
		model, cfg.MaxModelLen, a.TP, cfg.GpuMemoryUtilization,
	)
}

func BuildKube(model string, cfg ServingConfig, a Analysis) string {
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
    metadata:
      labels:
        app: pusula-serve
    spec:
      containers:
        - name: engine
          image: vllm/vllm-openai:latest
          args:
            - %s
            - --max-model-len
            - "%d"
            - --tensor-parallel-size
            - "%d"
            - --gpu-memory-utilization
            - "%.2f"
            - --max-num-seqs
            - "%d"
            - --enable-prefix-caching
          resources:
            limits:
              nvidia.com/gpu: %d
`, model, cfg.MaxModelLen, a.TP, cfg.GpuMemoryUtilization, cfg.MaxNumSeqs, cfg.NumGpus)
}
