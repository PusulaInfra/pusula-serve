package engine

import "math"

func round(v float64, digits int) float64 {
	if digits < 0 {
		digits = 0
	}
	p := math.Pow(10, float64(digits))
	return math.Round(v*p) / p
}
