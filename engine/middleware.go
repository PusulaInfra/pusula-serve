package engine

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func CardMiddleware(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		processor := Acquire(logger)
		defer Release(processor)

		ctx := context.WithValue(r.Context(), "cardProcessor", processor)
		r = r.WithContext(ctx)

		next(w, r)

		logger.Info("card request processed",
			"duration_ms", time.Since(start).Milliseconds(),
			"context_cuts", GlobalMetrics.ContextCuts.Load(),
		)
	}
}
