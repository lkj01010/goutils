package mathx

import (
	"math"
	"math/rand"
)

const KAccuracy float64 = 0.0000001

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

//func IsApproximate(a, b float64) bool {
func IsEqual(a, b float64) bool {
	return math.Abs(a-b) <= KAccuracy
}

func NotEqual(a, b float64) bool {
	return math.Abs(a-b) > KAccuracy
}

func IsBigger(a, b float64) bool {
	return a-b > KAccuracy
}

func IsSmaller(a, b float64) bool {
	return b-a > KAccuracy
}

// >=
func IsBiggerEqual(a, b float64) bool {
	return IsEqual(a, b) || IsBigger(a, b)
}

// <=
func IsSmallEqual(a, b float64) bool {
	return IsEqual(a, b) || IsSmaller(a, b)
}

func FloorInt(v float64) int {
	return int(math.Floor(v) + KAccuracy)
}

func CeilInt(v float64) int {
	return int(math.Ceil(v) + KAccuracy)
}

// std:lerp, unity use "t"
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func LerpIndex(fIndex float64) (int, int, float64) {
	left := math.Floor(fIndex)
	if IsEqual(left, fIndex) {
		return int(left), int(left + 1), 0
	} else {
		right := math.Ceil(fIndex)
		percent := (fIndex - left) / (right - left)
		return int(left), int(right), percent
	}
}

func RandUnit() float64 {
	r := rand.Intn(100)
	t := float64(r) / 100
	return t
}

func RandUnitWiggle() float64 {
	return RandUnit() - 0.5
}

func FloatInsect(a1, a2, b1, b2 float64) bool {
	if a1 > a2 {
		t := a1
		a1 = a2
		a2 = t
	}
	if b1 > b2 {
		t := b1
		b1 = b2
		b2 = t
	}
	if a1 >= b2 || a2 <= b1 {
		return false
	} else {
		return true
	}
}

func QuickDistance2D64(x, y float64) float64 {
	min := math.Min(x, y)
	max := math.Max(x, y)
	m := min / max
	q := m*m*0.37 + 0.05*m
	dist := max + max*q
	return dist
}
