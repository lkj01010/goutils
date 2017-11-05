package math

import "math"

var accuracy float64 = 0.0001

func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func ClampInt32(value, min, max int32) int32 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func EqualFloat64(a, b float64) bool {
	return math.Abs(a - b) < accuracy
}