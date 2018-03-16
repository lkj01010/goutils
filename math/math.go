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

func QuickDistance2D64(x, y float64) float64 {
	min := math.Min(x, y)
	max := math.Max(x, y)
	m := min / max
	q := m * m * 0.37 + 0.05 * m
	dist := max + max * q
	return dist
}