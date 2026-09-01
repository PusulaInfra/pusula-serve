package engine

type ModelRegistryEntry struct {
	Name            string  `json:"name"`
	TotalParamsB    float64 `json:"total_params_b"`
	ActiveParamsB   float64 `json:"active_params_b"`
	NLayers         int     `json:"n_layers"`
	NumKvHeads      int     `json:"num_kv_heads"`
	HiddenSize      int     `json:"hidden_size"`
	IsMoE           bool    `json:"is_moe"`
	NumExperts      int     `json:"num_experts"`
	UseMLALatent    bool    `json:"use_mla_latent"`
	MLALatentDim    int     `json:"mla_latent_dim"`
}

var GlobalModelRegistry = map[string]ModelRegistryEntry{
	"DeepSeek-R1": {
		Name:          "DeepSeek-R1",
		TotalParamsB:  671.0,
		ActiveParamsB: 37.0,
		NLayers:       61,
		NumKvHeads:    1, // MLA (Multi-head Latent Attention) sıkıştırılmış tek KV head yapısı
		HiddenSize:    7168,
		IsMoE:         true,
		NumExperts:    256,
		UseMLALatent:  true,
		MLALatentDim:  512,
	},
	"DeepSeek-V3": {
		Name:          "DeepSeek-V3",
		TotalParamsB:  671.0,
		ActiveParamsB: 37.0,
		NLayers:       61,
		NumKvHeads:    1,
		HiddenSize:    7168,
		IsMoE:         true,
		NumExperts:    256,
		UseMLALatent:  true,
		MLALatentDim:  512,
	},
	"Llama-3.1-70B": {
		Name:          "Llama-3.1-70B",
		TotalParamsB:  70.6,
		ActiveParamsB: 70.6,
		NLayers:       80,
		NumKvHeads:    8,
		HiddenSize:    8192,
		IsMoE:         false,
		NumExperts:    0,
		UseMLALatent:  false,
		MLALatentDim:  0,
	},
	"Qwen2.5-72B": {
		Name:          "Qwen2.5-72B",
		TotalParamsB:  72.7,
		ActiveParamsB: 72.7,
		NLayers:       80,
		NumKvHeads:    8,
		HiddenSize:    8192,
		IsMoE:         false,
		NumExperts:    0,
		UseMLALatent:  false,
		MLALatentDim:  0,
	},
}

// GetModelRegistryEntry, güncel model parametrelerini, katman sayılarını,
// MoE ve DeepSeek tarzı MLA (Multi-head Latent Attention) yapılandırmalarını getirir.
func GetModelRegistryEntry(modelName string) (ModelRegistryEntry, bool) {
	entry, exists := GlobalModelRegistry[modelName]
	if !exists {
		// Varsayılan standart Llama benzeri yapı
		return ModelRegistryEntry{
			Name:          modelName,
			TotalParamsB:  7.0,
			ActiveParamsB: 7.0,
			NLayers:       32,
			NumKvHeads:    8,
			HiddenSize:    4096,
			IsMoE:         false,
			NumExperts:    0,
			UseMLALatent:  false,
		}, false
	}
	return entry, true
}
