package httpapi

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
	Static fs.FS
}

type analyzeRequest struct {
	Model                string  `json:"model"`
	MaxModelLen          int     `json:"max_model_len"`
	NumGpus              int     `json:"num_gpus"`
	MaxNumSeqs           int     `json:"max_num_seqs"`
	DtypeBytes           int     `json:"dtype_bytes"`
	KVDtypeBytes         int     `json:"kv_dtype_bytes"`
	GpuType              string  `json:"gpu_type"`
	GpuMemoryUtilization float64 `json:"gpu_memory_utilization"`
	Engine               string  `json:"engine"`
	Provider             string  `json:"provider"`
	PrefixHit            float64 `json:"prefix_hit"`
}

func New() Server {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return Server{Static: sub}
}

func (s Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.health)
	mux.HandleFunc("/api/models", s.models)
	mux.HandleFunc("/api/analyze", s.analyze)
	mux.Handle("/", http.FileServer(http.FS(s.Static)))
	return mux
}

func (s Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "product": "pusula-serve"})
}

func (s Server) models(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"models": serve.ModelRegistry,
		"gpus":   serve.GPURegistry,
	})
}

func (s Server) analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	a := serve.Analyze(serve.ServingConfig{
		ModelName:            req.Model,
		MaxModelLen:          req.MaxModelLen,
		NumGpus:              req.NumGpus,
		MaxNumSeqs:           req.MaxNumSeqs,
		DtypeBytes:           req.DtypeBytes,
		KVDtypeBytes:         req.KVDtypeBytes,
		GpuType:              req.GpuType,
		GpuMemoryUtilization: req.GpuMemoryUtilization,
		Engine:               req.Engine,
		Provider:             req.Provider,
		PrefixHit:            req.PrefixHit,
	})
	writeJSON(w, http.StatusOK, a)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
