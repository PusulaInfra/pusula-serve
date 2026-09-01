package apply

import "testing"

func TestExecWithoutRemoteStaysDryRun(t *testing.T) {
	r, err := Run(Request{Mode: "exec", Remote: false, Line: "vllm serve x"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Mode != "dry-run" || r.Ran {
		t.Fatalf("exec without remote must dry-run, got %+v", r)
	}
}

func TestEmptyLineErrors(t *testing.T) {
	if _, err := Run(Request{}); err == nil {
		t.Fatal("expected error")
	}
}
