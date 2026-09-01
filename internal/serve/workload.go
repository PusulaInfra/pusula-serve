package serve

import "strings"

// WorkloadKind is the serving tax class. Not a robot OS.
// Each AI surface has a memory overhead; Pusula plans for it without running it.
type WorkloadKind string

const (
	WorkChat     WorkloadKind = "chat"
	WorkAgent    WorkloadKind = "agent"
	WorkEmbed    WorkloadKind = "embed"
	WorkRerank   WorkloadKind = "rerank"
	WorkWhisper  WorkloadKind = "whisper"
	WorkTTS      WorkloadKind = "tts"
	WorkVLM      WorkloadKind = "vlm"
	WorkSpec     WorkloadKind = "spec"
	WorkRAG      WorkloadKind = "rag"
	WorkEdge     WorkloadKind = "edge"
)

// Workload describes the serving task tax, not the platform.
// Every AI surface has a serving cost class. Pusula does not run them,
// it tells you if your box holds them.
type Workload struct {
	Kind         WorkloadKind `json:"kind"`
	Tools        int          `json:"tools"`          // agent tool schemas
	ImagesPerReq int          `json:"images_per_req"` // VLM
	ImageTokens  int          `json:"image_tokens"`   // tokens/image after tiles
	LoRAAdapters int          `json:"lora_adapters"`
	LoRAGBEach   float64      `json:"lora_gb_each"`
	DraftParamsB float64      `json:"draft_params_b"` // spec decode
	EmbedBatch   int          `json:"embed_batch"`
	SidecarGB    float64      `json:"sidecar_gb"` // rag embedder / whisper encoder
	AgentLoops   int          `json:"agent_loops"`
}

// ParseKind maps user input to WorkloadKind.
func ParseKind(s string) WorkloadKind {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "agent", "tools", "assistant":
		return WorkAgent
	case "embed", "embedding":
		return WorkEmbed
	case "rerank":
		return WorkRerank
	case "whisper", "asr", "stt":
		return WorkWhisper
	case "tts", "speech":
		return WorkTTS
	case "vlm", "vision", "multimodal":
		return WorkVLM
	case "spec", "speculative":
		return WorkSpec
	case "rag":
		return WorkRAG
	case "edge", "robot", "jetson":
		return WorkEdge
	default:
		return WorkChat
	}
}

// TaxGB returns extra GB beyond base weights + KV runtime.
// Never invents tok/s. Each workload class has its own memory overhead.
func (w Workload) TaxGB(weightGB float64, kvPerTokKB float64, ctx int, seqs int) (extraGB float64, note string) {
	if w.Kind == "" {
		w.Kind = WorkChat
	}

	// LoRA adapters: typically 0.35 GB each if not specified
	lora := float64(w.LoRAAdapters) * w.LoRAGBEach
	if w.LoRAGBEach == 0 && w.LoRAAdapters > 0 {
		lora = float64(w.LoRAAdapters) * 0.35
	}

	// Speculative decode draft model
	draft := 0.0
	if w.DraftParamsB > 0 {
		if w.DraftParamsB < 100 {
			draft = w.DraftParamsB * 2.0 // params_b * 2 bytes (BF16)
		} else {
			draft = w.DraftParamsB * 2.0 / 1024.0 * 2 // rough adjustment for large numbers
		}
	}

	sidecar := w.SidecarGB

	// VLM image token tax
	images := 0.0
	if w.ImageTokens > 0 {
		images = float64(w.ImagesPerReq*w.ImageTokens*seqs) * kvPerTokKB / (1024.0 * 1024.0)
	} else if w.ImagesPerReq > 0 {
		images = float64(w.ImagesPerReq*seqs) * 0.15 // ~tile tax fallback
	}

	// Agent tool & loop tax
	tools := 0.0
	if w.Kind == WorkAgent {
		tools = float64(w.Tools) * 0.02 // tool schemas are tiny
		if w.AgentLoops > 1 {
			// extra seq residency for tool loop KV
			toolKV := float64(w.AgentLoops-1) * kvPerTokKB * float64(ctx) * float64(seqs) * 0.15 / (1024.0 * 1024.0)
			tools += toolKV
		}
	}

	switch w.Kind {
	case WorkEmbed, WorkRerank:
		// No KV residency; batch activations instead
		act := float64(max(w.EmbedBatch, seqs)) * 0.08
		return lora + sidecar + act, "embed/rerank: KV lane off, batch activations on"

	case WorkWhisper, WorkTTS:
		// Encoder/decoder workspace, no chat KV
		return sidecar + 1.5 + lora, "audio: no chat KV; encoder/decoder workspace"

	case WorkVLM:
		return lora + sidecar + images, "vlm: image token tax on KV path"

	case WorkSpec:
		// Draft model is the tax, NOT a promise of 2x tok/s
		return lora + draft + sidecar, "spec: draft weights are the tax, not free tok/s"

	case WorkRAG:
		if sidecar == 0 {
			sidecar = 2.0 // default embed sidecar
		}
		return lora + sidecar + tools, "rag: embed sidecar + prefix-hit is the lever"

	case WorkEdge:
		// UMA ≠ H100. Jetson / Spark
		return lora + sidecar, "edge/robot: UMA ≠ 80GB H100"

	case WorkAgent:
		// Loops and tools tax residency. Not a robot OS.
		return lora + sidecar + tools, "agent: loops and tools tax residency"

	default:
		// WorkChat
		return lora + sidecar, "chat"
	}
}

// APIShape returns the OpenAI-compatible endpoint shape for this workload.
// Pusula does not run the API; it plans for it.
func APIShape(kind WorkloadKind) string {
	switch kind {
	case WorkEmbed:
		return "/v1/embeddings"
	case WorkRerank:
		return "/v1/rerank"
	case WorkWhisper:
		return "/v1/audio/transcriptions"
	case WorkTTS:
		return "/v1/audio/speech"
	case WorkVLM, WorkAgent, WorkChat, WorkRAG, WorkSpec:
		return "/v1/chat/completions"
	case WorkEdge:
		return "/v1/chat/completions (edge, batch=1)"
	default:
		return "/v1/chat/completions"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
