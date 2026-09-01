package engine

import (
	"encoding/json"
	"fmt"
)

type DeploymentConfig struct {
	ModelName  string `json:"model_name"`
	Replicas   int    `json:"replicas"`
	TotalGPUs  int    `json:"total_gpus"`
	GPUFlavor  string `json:"gpu_flavor"`
	MaxContext int    `json:"max_context"`
}

// GenerateExports, yapılan cluster hesaplama konfigürasyonunu 
// JSON, YAML, Kubernetes Helm değerleri (values.yaml) ve Terraform snippet formatına dönüştürür.
func GenerateExports(cfg DeploymentConfig) map[string]string {
	// 1. JSON Export
	jsonBytes, _ := json.MarshalIndent(cfg, "", "  ")
	jsonStr := string(jsonBytes)

	// 2. YAML Export
	yamlStr := fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: pusula-serve-config
data:
  model_name: "%s"
  replicas: "%d"
  total_gpus: "%d"
  gpu_flavor: "%s"
  max_context: "%d"`, cfg.ModelName, cfg.Replicas, cfg.TotalGPUs, cfg.GPUFlavor, cfg.MaxContext)

	// 3. Helm Values Snippet
	helmStr := fmt.Sprintf(`replicaCount: %d
image:
  repository: vllm/vllm-openai
  tag: latest
resources:
  limits:
    nvidia.com/gpu: %d
env:
  - name: MODEL_NAME
    value: "%s"
  - name: MAX_MODEL_LEN
    value: "%d"`, cfg.Replicas, cfg.TotalGPUs, cfg.ModelName, cfg.MaxContext)

	// 4. Terraform Snippet
	terraformStr := fmt.Sprintf(`resource "kubernetes_deployment" "llm_serve" {
  metadata {
    name = "pusula-%s"
  }
  spec {
    replicas = %d
    # GPU Flavor: %s, Max Context: %d
  }
}`, cfg.ModelName, cfg.Replicas, cfg.GPUFlavor, cfg.MaxContext)

	return map[string]string{
		"json":      jsonStr,
		"yaml":      yamlStr,
		"helm":      helmStr,
		"terraform": terraformStr,
	}
}
