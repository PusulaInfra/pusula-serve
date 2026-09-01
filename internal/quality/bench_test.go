package quality

import "testing"

func TestRunDoesNotInventThroughput(t *testing.T) {
	r := Run(7)
	if r.Ran || r.TokPerS != 0 {
		t.Fatalf("must not invent tok/s: %+v", r)
	}
	if r.Seconds != 7 {
		t.Fatalf("seconds=%d", r.Seconds)
	}
}
