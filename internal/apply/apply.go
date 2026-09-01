package apply

import "errors"

type Request struct {
	Mode   string `json:"mode"`   // dry-run | exec
	Remote bool   `json:"remote"`
	Line   string `json:"line"`
}

type Result struct {
	Mode    string `json:"mode"`
	Ran     bool   `json:"ran"`
	Line    string `json:"line"`
	Note    string `json:"note"`
}

func Run(req Request) (Result, error) {
	if req.Line == "" {
		return Result{}, errors.New("empty launch line")
	}
	mode := req.Mode
	if mode == "" {
		mode = "dry-run"
	}
	if mode == "exec" && !req.Remote {
		return Result{Mode: "dry-run", Ran: false, Line: req.Line, Note: "exec refused without --remote"}, nil
	}
	if mode != "exec" {
		return Result{Mode: "dry-run", Ran: false, Line: req.Line, Note: "printed only; not started"}, nil
	}
	return Result{Mode: "exec", Ran: false, Line: req.Line, Note: "exec acknowledged; this process does not spawn engines"}, nil
}
