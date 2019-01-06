package r2

import (
	"fmt"
	"math"
	mymath "github.com/lkj01010/goutils/math"
)

type Vec struct {
	X, Y float64
}

// Add returns the sum of v and ov.
func (v Vec) Add(ov Vec) Vec { return Vec{v.X + ov.X, v.Y + ov.Y} }

// Sub returns the difference of v and ov.
func (v Vec) Sub(ov Vec) Vec { return Vec{v.X - ov.X, v.Y - ov.Y} }

// Mul returns the scalar product of v and m.
func (v Vec) Mul(m float64) Vec { return Vec{m * v.X, m * v.Y} }

// Ortho returns a counterclockwise orthogonal point with the same norm.
func (v Vec) Ortho() Vec { return Vec{-v.Y, v.X} }

// Dot returns the dot product between v and ov.
func (v Vec) Dot(ov Vec) float64 { return v.X*ov.X + v.Y*ov.Y }

// Cross returns the cross product of v and ov.
func (v Vec) Cross(ov Vec) float64 { return v.X*ov.Y - v.Y*ov.X }

// Norm returns the vector's Length.
func (v Vec) Length() float64 { return math.Hypot(v.X, v.Y) }

// Normalize returns a unit point in the same direction as v.
func (v Vec) Normalize() Vec {
	if v.X == 0 && v.Y == 0 {
		return v
	}
	return v.Mul(1 / v.Length())
}

func (v Vec) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vec) Scale(xf, yf float64) {
	v.X *= xf
	v.Y *= yf
}

func (v Vec) Equals(ov Vec) bool {
	return mymath.EqualFloat64(v.X, ov.X) && mymath.EqualFloat64(v.Y, ov.Y)
	//return v.X == ov.X && v.Y == ov.Y
}

func (v Vec) String() string { return fmt.Sprintf("(%.12f, %.12f)", v.X, v.Y) }

func (v Vec) AngleDeg() float64 {
	angle := math.Atan2(v.Y, v.X) * mymath.Rad2Deg
	if angle < 0 {
		angle += 360
	}
	return angle
}

func (v Vec) AngleRad() float64 {
	return math.Atan2(v.Y, v.X)
}


func NewVecDirFromAngleDeg(angleDeg float64) Vec {
	x := math.Cos(angleDeg * mymath.Deg2Rad)
	y := math.Sin(angleDeg * mymath.Deg2Rad)
	return Vec{x, y}
}

//--------------------------------------------------

type Rot struct {
	S, C float64
}

func (r *Rot) Set(angleRad float64) {
	r.S = math.Sin(angleRad)
	r.C = math.Cos(angleRad)
}

func (r Rot) GetAngleRad() float64 {
	return math.Atan2(r.S, r.C)
}

//--------------------------------------------------

type Transform struct {
	P Vec
	Q Rot
}

func (t *Transform) Set(pos *Vec, angleRad float64) {
	t.P = *pos
	t.Q.Set(angleRad)
}
