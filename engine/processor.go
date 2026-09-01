package engine

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	MaxSeqStands = 16
	MaxSeqPages  = 32
)

var (
	ErrMaxStandsExceeded = errors.New("limit exceeded: maximum 16 sequence stands allowed")
	ErrMaxPagesExceeded  = errors.New("limit exceeded: maximum 32 sequence pages allowed")
)

type Metrics struct {
	ActiveStands  atomic.Int64
	ActivePages   atomic.Int64
	ContextCuts   atomic.Int64
	TotalErrors   atomic.Int64
}

var GlobalMetrics = &Metrics{}

type CardProcessor struct {
	mu        sync.RWMutex
	seqStands int
	seqPages  int
	context   map[string]any
	logger    *slog.Logger
}

var ProcessorPool = sync.Pool{
	New: func() any {
		return &CardProcessor{
			context: make(map[string]any, 16),
			logger:  slog.Default(),
		}
	},
}

func Acquire(logger *slog.Logger) *CardProcessor {
	cp := ProcessorPool.Get().(*CardProcessor)
	cp.mu.Lock()
	if logger != nil {
		cp.logger = logger
	}
	cp.mu.Unlock()
	return cp
}

func Release(cp *CardProcessor) {
	cp.CutContext(context.Background())
	ProcessorPool.Put(cp)
}

func (cp *CardProcessor) AddPage(ctx context.Context) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.seqPages >= MaxSeqPages {
		GlobalMetrics.TotalErrors.Add(1)
		cp.logger.WarnContext(ctx, "page limit reached", "pages", cp.seqPages)
		return ErrMaxPagesExceeded
	}
	cp.seqPages++
	GlobalMetrics.ActivePages.Store(int64(cp.seqPages))
	return nil
}

func (cp *CardProcessor) AddSequenceStand(ctx context.Context) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.seqStands >= MaxSeqStands {
		GlobalMetrics.TotalErrors.Add(1)
		cp.logger.WarnContext(ctx, "sequence stand limit reached", "stands", cp.seqStands)
		return ErrMaxStandsExceeded
	}
	cp.seqStands++
	cp.seqPages = 0
	GlobalMetrics.ActiveStands.Store(int64(cp.seqStands))
	GlobalMetrics.ActivePages.Store(0)
	return nil
}

func (cp *CardProcessor) CutContext(ctx context.Context) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for k := range cp.context {
		delete(cp.context, k)
	}
	cp.seqStands = 0
	cp.seqPages = 0
	
	GlobalMetrics.ContextCuts.Add(1)
	GlobalMetrics.ActiveStands.Store(0)
	GlobalMetrics.ActivePages.Store(0)

	cp.logger.DebugContext(ctx, "card context cut successfully executed")
}

func (cp *CardProcessor) SetContextValue(key string, value any) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.context[key] = value
}

func (cp *CardProcessor) GetContextValue(key string) (any, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	val, ok := cp.context[key]
	return val, ok
}
