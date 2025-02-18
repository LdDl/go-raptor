package raptor

import "math"

// MinInt returns the minimum of two integers.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Infinity is the maximum value for an 32-bit integer.
const Infinity = math.MaxInt32
