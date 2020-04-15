package primitives

import (
	"math"
)

type Triangle struct {
	N          Vector
	P1, P2, P3 Point
}

func NewTriangle(p1, p2, p3 Point) Triangle {
	t := Triangle{}
	t.Fill(p1, p2, p3)
	return t
}

func (t *Triangle) Shift(v Vector) Triangle {
	return NewTriangle(t.P1.Shift(v), t.P2.Shift(v), t.P3.Shift(v))
}

func (t *Triangle) Fill(p1, p2, p3 Point) {
	t.P1 = p1
	t.P2 = p2
	t.P3 = p3
	t.recalculate()
}

func (t *Triangle) recalculate() {
	t.N = normal(t.P1.VectorTo(t.P2), t.P1.VectorTo(t.P3))
}

func normal(v1, v2 Vector) Vector {
	return v1.Cross(v2)
}

func (t *Triangle) MinZ(z Vector) float64 { // it is without normalization!!
	pr1 := t.P1.ToVector().Dot(z)
	pr2 := t.P2.ToVector().Dot(z)
	pr3 := t.P3.ToVector().Dot(z)
	return math.Min(pr1, math.Min(pr2, pr3))
}

func (t *Triangle) MaxZ(z Vector) float64 {
	pr1 := t.P1.ToVector().Dot(z)
	pr2 := t.P2.ToVector().Dot(z)
	pr3 := t.P3.ToVector().Dot(z)
	return math.Max(pr1, math.Max(pr2, pr3))
}

func (t *Triangle) MinMaxZ(z Vector) (float64, float64) {
	pr1 := t.P1.ToVector().Dot(z)
	pr2 := t.P2.ToVector().Dot(z)
	pr3 := t.P3.ToVector().Dot(z)
	return math.Min(pr1, math.Min(pr2, pr3)), math.Max(pr1, math.Max(pr2, pr3))
}

func (t *Triangle)IntersectTriangle(t2 *Triangle) *Line {

	plane := Plane{P: t.P1, N: t.N}
	line := plane.IntersectTriangle(t2)
	if line == nil{
		return nil
	}

	b1 := t.PointBelongs(line.P1)
	b2 := t.PointBelongs(line.P2)

	//if line is inside of triangle
	if b1 && b2 {
		return line
	}

	p1 := Line{P1: t.P1, P2: t.P2}.IntersectLine(line)
	p2 := Line{P1: t.P1, P2: t.P3}.IntersectLine(line)
	p3 := Line{P1: t.P2, P2: t.P3}.IntersectLine(line)

	//if line crosses the one side of triangle
	if b1 || b2 {
		if p1 != nil{
			if b1{
				return &Line{P1: line.P1, P2: *p1}
			}
			if b2{
				return &Line{P1: *p1, P2: line.P2}
			}
		}
		if p2 != nil{
			if b1{
				return &Line{P1: line.P1, P2: *p2}
			}
			if b2{
				return &Line{P1: *p2, P2: line.P2}
			}
		}
		if p3 != nil{
			if b1{
				return &Line{P1: line.P1, P2: *p3}
			}
			if b2{
				return &Line{P1: *p3, P2: line.P2}
			}
		}

	}

	// if line crosses 2 sides of triangle
	if p1 == nil && p2 != nil && p3 != nil{
		return &Line{P1: *p2, P2: *p3}
	}

	if p1 != nil && p2 == nil && p3 != nil{
		return &Line{P1: *p1, P2: *p3}
	}

	if p1 != nil && p2 != nil && p3 == nil{
		return &Line{P1: *p1, P2: *p2}
	}

	// if line cross all siides of triangle (one edge and one vertex)
	if p1 != nil && p2 != nil && p3 != nil {
		if p1 == p2 {
			return &Line{P1: *p1, P2: *p3}
		}
		if p1  == p3 {
			return &Line{P1: *p1, P2: *p2}
		}
		if p2  == p3 {
			return &Line{P1: *p1, P2: *p2}
		}
	}

	return nil
}

func (t *Triangle) PointBelongs(p Point) bool {

	plane := Plane{P: t.P1, N: t.N}
	if plane.PointBelongs(p) {

		A := t.P1
		B := t.P2
		C := t.P3

		w1 := (A.X - p.X)*(B.Y - A.Y) - (B.X - A.X)*(A.Y - p.Y)
		w2 := (B.X - p.X)*(C.Y - B.Y) - (C.X - B.X)*(B.Y - p.Y)
		w3 := (C.X - p.X)*(A.Y - C.Y) - (A.X - C.X)*(C.Y - p.Y)

		if w1>=0 && w2>=0 && w3 >=0 || w1<=0 && w2<=0 && w3 <=0 {
			return true
		}
	}

	return false
}