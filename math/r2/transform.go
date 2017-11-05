package r2

import "math"
import mymath "github.com/lkj01010/goutils/math"

/// Rotate a vector
func VecRot(q Rot, v Vec) Vec {
	return Vec{
		q.C*v.X - q.S*v.Y,
		q.S*v.X + q.C*v.Y,
	}
}

func VecRotInverse(q Rot, v Vec) Vec {
	return Vec{
		q.C*v.X + q.S*v.Y,
		- q.S*v.X + q.C*v.Y,
	}
}

func VecTransform(t Transform, v Vec) Vec {
	x := (t.Q.C*v.X - t.Q.S*v.Y) + t.P.X;
	y := (t.Q.S*v.X + t.Q.C*v.Y) + t.P.Y;
	return Vec{x, y}
}

func VecTransformInverse(t Transform, v Vec) Vec {
	px := v.X - t.P.X;
	py := v.Y - t.P.Y;
	x := t.Q.C*px + t.Q.S*py;
	y := -t.Q.S*px + t.Q.C*py;
	return Vec{x, y}
}

func RotRot(q, r Rot) (qr Rot) {
	qr.S = q.S*r.C + q.C*r.S
	qr.C = q.C*r.C - q.S*r.S
	return
}

func RotRotInverse(q, r Rot) (qr Rot) {
	qr.S = q.C*r.S - q.S*r.C
	qr.C = q.C*r.C + q.S*r.S
	return
}

func TransformTransform(a Transform, b Transform) (c Transform) {
	c.Q = RotRot(a.Q, b.Q)
	c.P = VecRot(a.Q, b.P).Add(a.P)
	return
}

func TransformTransformInverse(a Transform, b Transform) (c Transform) {
	c.Q = RotRotInverse(a.Q, b.Q)
	c.P = VecRotInverse(a.Q, b.P.Sub(a.P))
	return
}

func AngleBetweenVec(from Vec, to Vec) float64 {
	cosValue := from.Dot(to) / (from.Length() * to.Length())
	cosValue = mymath.ClampFloat64(cosValue, -1, 1)
	return math.Acos(cosValue)
}
