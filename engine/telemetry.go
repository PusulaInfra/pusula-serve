package engine

import (
	"encoding/json"
	"net/http"
)

// TelemetryData represents the structural state of the engine metrics.
type TelemetryData struct {
	ActiveStands int64 `json:"active_stands"`
	ActivePages  int64 `json:"active_pages"`
	ContextCuts  int64 `json:"context_cuts"`
	TotalErrors  int64 `json:"total_errors"`
}

// RegisterTelemetryRoutes mounts monitoring and health check endpoints onto a ServeMux.
func RegisterTelemetryRoutes(mux *http.ServeMux) {
	// Prometheus / JSON Metrics Endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		data := TelemetryData{
			ActiveStands: GlobalMetrics.ActiveStands.Load(),
			ActivePages:  GlobalMetrics.ActivePages.Load(),
			ContextCuts:  GlobalMetrics.ContextCuts.Load(),
			TotalErrors:  GlobalMetrics.TotalErrors.Load(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(data)
	})

	// Kubernetes Liveness Probe
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Kubernetes Readiness Probe
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Add custom checks here if model weights or engine pools are warming up
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY"))
	})
}
