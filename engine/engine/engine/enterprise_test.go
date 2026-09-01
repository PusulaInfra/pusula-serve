package engine

import (
	"context"
	"sync"
	"testing"
)

func TestEngineStressAndQueue(t *testing.T) {
	ctx := context.Background()
	var wg sync.WaitGroup
	workers := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := GlobalQueue.AcquireSlot(ctx)
			if err == nil {
				defer GlobalQueue.ReleaseSlot()
				cp := Acquire(nil)
				defer Release(cp)
				_ = cp.AddSequenceStand(ctx)
				_ = cp.AddPage(ctx)
			}
		}()
	}
	wg.Wait()
}
