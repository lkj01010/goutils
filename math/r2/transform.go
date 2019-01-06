package r2

import "math"
import mymath "github.com/lkj01010/goutils/math"

/// Rotate a vector
/*
x = Rcos(b) ; y = Rsin(b);
X = Rcos(a+b) = Rcosacosb - Rsinasinb = xcosa - ysina; (合角公式)
Y = Rsin(a+b) = Rsinacosb + Rcosasinb = xsina + ycosa ;
*/
func VecRot(v Vec, q Rot) Vec {
	return Vec{
		v.X*q.C - v.Y*q.S,
		v.X*q.S + v.Y*q.C,
	}
}

func VecRotInverse(v Vec, q Rot) Vec {
	return Vec{
		v.X*q.C + v.Y*q.S,
		- v.X*q.S + v.Y*q.C,
	}
}

func VecTransform(v Vec, t Transform) Vec {
	x := (t.Q.C*v.X - t.Q.S*v.Y) + t.P.X;
	y := (t.Q.S*v.X + t.Q.C*v.Y) + t.P.Y;
	return Vec{x, y}
}

func VecTransformInverse(v Vec, t Transform) Vec {
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
	c.P = VecRot(b.P, a.Q).Add(a.P)
	return
}

func TransformTransformInverse(a Transform, b Transform) (c Transform) {
	c.Q = RotRotInverse(a.Q, b.Q)
	c.P = VecRotInverse(b.P.Sub(a.P), a.Q)
	return
}

func AngleBetweenVec(from Vec, to Vec) float64 {
	cosValue := from.Dot(to) / (from.Length() * to.Length())
	cosValue = mymath.ClampFloat64(cosValue, -1, 1)
	return math.Acos(cosValue)
}
