package engine

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// --- 1. Distributed Tracing & Telemetry ---
type TraceMetric struct {
	TTFTMillis      atomic.Int64
	TokenGenSpeed   atomic.Int64
	ActiveSpanCount atomic.Int64
}

var GlobalTrace = &TraceMetric{}

type TraceSpan struct {
	Name      string
	StartTime time.Time
}

func StartSpan(name string) *TraceSpan {
	GlobalTrace.ActiveSpanCount.Add(1)
	return &TraceSpan{
		Name:      name,
		StartTime: time.Now(),
	}
}

func (s *TraceSpan) End(ctx context.Context, ttft time.Duration, tokensPerSec int64) {
	GlobalTrace.ActiveSpanCount.Add(-1)
	GlobalTrace.TTFTMillis.Store(ttft.Milliseconds())
	GlobalTrace.TokenGenSpeed.Store(tokensPerSec)
}

// --- 2. Smart Queueing & Rate Limiting ---
var ErrQueueTimeout = errors.New("request queue timeout: engine capacity fully saturated")

type RequestQueue struct {
	semaphore chan struct{}
	timeout   time.Duration
}

var GlobalQueue = NewRequestQueue(16, 2*time.Second)

func NewRequestQueue(maxConcurrent int, timeout time.Duration) *RequestQueue {
	return &RequestQueue{
		semaphore: make(chan struct{}, maxConcurrent),
		timeout:   timeout,
	}
}

func (q *RequestQueue) AcquireSlot(ctx context.Context) error {
	select {
	case q.semaphore <- struct{}{}:
		return nil
	case <-time.After(q.timeout):
		GlobalMetrics.TotalErrors.Add(1)
		return ErrQueueTimeout
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *RequestQueue) ReleaseSlot() {
	select {
	case <-q.semaphore:
	default:
	}
}

// --- 3. Hot-Reload Configuration Management ---
type EngineConfig struct {
	ModelName   string
	MaxModelLen int
	MaxSeqs     int
}

type ConfigManager struct {
	currentConfig atomic.Value
}

var GlobalConfigManager = &ConfigManager{}

func init() {
	GlobalConfigManager.currentConfig.Store(EngineConfig{
		ModelName:   "llama-3.1-70b",
		MaxModelLen: 16384,
		MaxSeqs:     16,
	})
}

func (cm *ConfigManager) Get() EngineConfig {
	return cm.currentConfig.Load().(EngineConfig)
}

func (cm *ConfigManager) Update(newCfg EngineConfig) {
	cm.currentConfig.Store(newCfg)
}
