package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/PusulaInfra/pusula-serve/internal/apply"
	"github.com/PusulaInfra/pusula-serve/internal/live"
	"github.com/PusulaInfra/pusula-serve/internal/measure"
	"github.com/PusulaInfra/pusula-serve/internal/serve"
)

type Server struct{}

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
	BenchSec             int     `json:"bench_sec"`
	ApplyMode            string  `json:"apply"`
	Remote               bool    `json:"remote"`
	Line                 string  `json:"line"`
}

func New() Server { return Server{} }

func (s Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.health)
	mux.HandleFunc("/api/models", s.models)
	mux.HandleFunc("/api/analyze", s.analyze)
	mux.HandleFunc("/api/measure", s.measure)
	mux.HandleFunc("/api/live-vram", s.liveVRAM)
	mux.HandleFunc("/api/apply", s.apply)
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

func decode(w http.ResponseWriter, r *http.Request) (analyzeRequest, bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return analyzeRequest{}, false
	}
	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return analyzeRequest{}, false
	}
	return req, true
}

func cfgOf(req analyzeRequest) serve.ServingConfig {
	return serve.ServingConfig{
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
	}
}

func (s Server) analyze(w http.ResponseWriter, r *http.Request) {
	req, ok := decode(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, serve.Analyze(cfgOf(req)))
}

func (s Server) measure(w http.ResponseWriter, r *http.Request) {
	req, ok := decode(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, measure.Run(cfgOf(req), req.BenchSec))
}

func (s Server) liveVRAM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, live.SnapshotNow())
}

func (s Server) apply(w http.ResponseWriter, r *http.Request) {
	req, ok := decode(w, r)
	if !ok {
		return
	}
	line := req.Line
	if line == "" {
		line = serve.Analyze(cfgOf(req)).EngineCommand
	}
	res, err := apply.Run(apply.Request{Mode: req.ApplyMode, Remote: req.Remote, Line: line})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
