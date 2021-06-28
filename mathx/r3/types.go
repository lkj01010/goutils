package r3

import (
    "fmt"
    "github.com/lkj01010/goutils/mathx"
    "math"
)

type Vec struct {
    X, Y, Z float64
}

// Add returns the sum of v and ov.
func (v Vec) Add(ov Vec) Vec { return Vec{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z} }

// Sub returns the difference of v and ov.
func (v Vec) Sub(ov Vec) Vec { return Vec{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z} }

// Mul returns the scalar product of v and m.
func (v Vec) Mul(m float64) Vec { return Vec{m * v.X, m * v.Y, m * v.Z} }

// Ortho returns a counterclockwise orthogonal point with the same norm.
//func (v Vec) Ortho() Vec { return Vec{-v.Y, v.X} }

// Dot returns the dot product between v and ov.
func (v Vec) Dot(ov Vec) float64 { return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z }

// Cross returns the cross product of v and ov.
//func (v Vec) Cross(ov Vec) float64 { return v.X*ov.Y - v.Y*ov.X }

// Norm returns the vector's Length.
func (v Vec) Length() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z) }

// Normalize returns a unit point in the same direction as v.
func (v Vec) Normalize() Vec {
    if v.X == 0 && v.Y == 0 && v.Z == 0 {
        return v
    }
    return v.Mul(1 / v.Length())
}

func (v Vec) LengthSquared() float64 {
    return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec) Scale(xf, yf, zf float64) {
    v.X *= xf
    v.Y *= yf
    v.Z *= zf
}

func (v Vec) Equals(ov Vec) bool {
    return mathx.IsEqual(v.X, ov.X) && mathx.IsEqual(v.Y, ov.Y) && mathx.IsEqual(v.Z, ov.Z)
    //return v.X == ov.X && v.Y == ov.Y
}

func (v Vec) String() string { return fmt.Sprintf("(%.12f, %.12f, %.12f)", v.X, v.Y, v.Z) }

//func (v Vec) AngleDeg() float64 {
//    angle := math.Atan2(v.Y, v.X) * mathx.Rad2Deg
//    if angle < 0 {
//        angle += 360
//    }
//    return angle
//}
//
//func (v Vec) AngleRad() float64 {
//    return math.Atan2(v.Y, v.X)
//}
//
//func NewVecDirFromAngleDeg(angleDeg float64) Vec {
//    x := math.Cos(angleDeg * mathx.Deg2Rad)
//    y := math.Sin(angleDeg * mathx.Deg2Rad)
//    return Vec{x, y}
//}

//--------------------------------------------------
