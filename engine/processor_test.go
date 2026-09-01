package engine

import (
	"context"
	"testing"
)

func TestCardProcessorLimits(t *testing.T) {
	cp := Acquire(nil)
	defer Release(cp)
	ctx := context.Background()

	for i := 0; i < MaxSeqStands; i++ {
		if err := cp.AddSequenceStand(ctx); err != nil {
			t.Fatalf("failed at stand %d: %v", i, err)
		}
	}

	if err := cp.AddSequenceStand(ctx); err != ErrMaxStandsExceeded {
		t.Errorf("expected ErrMaxStandsExceeded, got %v", err)
	}

	for i := 0; i < MaxSeqPages; i++ {
		if err := cp.AddPage(ctx); err != nil {
			t.Fatalf("failed at page %d: %v", i, err)
		}
	}

	if err := cp.AddPage(ctx); err != ErrMaxPagesExceeded {
		t.Errorf("expected ErrMaxPagesExceeded, got %v", err)
	}

	cp.CutContext(ctx)
	if cp.seqStands != 0 || cp.seqPages != 0 {
		t.Errorf("CutContext failed to reset counters")
	}
}
