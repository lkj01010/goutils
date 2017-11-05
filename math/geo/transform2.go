package geo

import "math"

type Transform2 struct {
	p Vec2
	q Rot2
}

func Translate(v *Vec2, d Vec2) {
	v.Add(d)
}

func Rotate(v *Vec2, angle float32) {

}

//////////////////////////////////////////////////

type Rot2 struct {
	s, c float64
}

func (r *Rot2) Set(angle float64) {
	r.s = math.Sin(angle)
	r.c = math.Cos(angle)
}

/// Multiply two rotations: q * r
func (r *Rot2) Mul(or *Rot2) {
	// [qc -qs] * [rc -rs] = [qc*rc-qs*rs -qc*rs-qs*rc]
	// [qs  qc]   [rs  rc]   [qs*rc+qc*rs -qs*rs+qc*rc]
	// s = qs * rc + qc * rs
	// c = qc * rc - qs * rs
	var s, c float64
	s = r.s * or.c + r.c * or.s
	c = r.c * or.c - r.s * or.s
	r.s, r.c = s, c
}

/// Transpose multiply two rotations: qT * r
func (r *Rot2) MulT(or *Rot2) {
	// [ qc qs] * [rc -rs] = [qc*rc+qs*rs -qc*rs+qs*rc]
	// [-qs qc]   [rs  rc]   [-qs*rc+qc*rs qs*rs+qc*rc]
	// s = qc * rs - qs * rc
	// c = qc * rc + qs * rs
}

func (r Rot2) Apply(v *Vec2) {

}


