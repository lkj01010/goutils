package geo

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X, Y float64
}

// Add returns the sum of v and ov.
func (v Vec2) Add(ov Vec2) Vec2 { return Vec2{v.X + ov.X, v.Y + ov.Y} }

// Sub returns the difference of v and ov.
func (v Vec2) Sub(ov Vec2) Vec2 { return Vec2{v.X - ov.X, v.Y - ov.Y} }

// Mul returns the scalar product of v and m.
func (v Vec2) Mul(m float64) Vec2 { return Vec2{m * v.X, m * v.Y} }

// Ortho returns a counterclockwise orthogonal point with the same norm.
func (v Vec2) Ortho() Vec2 { return Vec2{-v.Y, v.X} }

// Dot returns the dot product between v and ov.
func (v Vec2) Dot(ov Vec2) float64 { return v.X*ov.X + v.Y*ov.Y }

// Cross returns the cross product of v and ov.
func (v Vec2) Cross(ov Vec2) float64 { return v.X*ov.Y - v.Y*ov.X }

// Norm returns the vector's norm.
func (v Vec2) Norm() float64 { return math.Hypot(v.X, v.Y) }

// Normalize returns a unit point in the same direction as v.
func (v Vec2) Normalize() Vec2 {
	if v.X == 0 && v.Y == 0 {
		return v
	}
	return v.Mul(1 / v.Norm())
}

func (v *Vec2) Scale(xf, yf float64) {
	v.X *= xf
	v.Y *= yf
}

func (v Vec2) Equals(ov Vec2) bool { return v.X == ov.X && v.Y == ov.Y }

func (v Vec2) String() string { return fmt.Sprintf("(%.12f, %.12f)", v.X, v.Y) }
