package quality

const Disclaimer = "Measured on this box, this run. Not a vendor SLA."

type Result struct {
	Ran        bool    `json:"ran"`
	Seconds    int     `json:"seconds"`
	TokPerS    float64 `json:"tok_per_s"`
	Disclaimer string  `json:"disclaimer"`
	Note       string  `json:"note"`
}

// Run does not call an engine. A missing engine is a skip, not a number.
func Run(seconds int) Result {
	if seconds < 1 {
		seconds = 5
	}
	return Result{
		Ran:        false,
		Seconds:    seconds,
		TokPerS:    0,
		Disclaimer: Disclaimer,
		Note:       "no local engine hooked; tok/s not invented",
	}
}
