package goosli

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

func (t *Triangle) MinZSquareDirected(z Vector) float64 {
	pr1 := t.P1.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr2 := t.P2.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr3 := t.P3.ToVector().ProjectOn(z).LengthSquareDirected(z)
	return math.Min(pr1, math.Min(pr2, pr3))
}

func (t *Triangle) MaxZSquareDirected(z Vector) float64 {
	pr1 := t.P1.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr2 := t.P2.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr3 := t.P3.ToVector().ProjectOn(z).LengthSquareDirected(z)
	return math.Max(pr1, math.Max(pr2, pr3))
}

func (t *Triangle) MinMaxZSquareDirected(z Vector) (float64, float64) {
	pr1 := t.P1.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr2 := t.P2.ToVector().ProjectOn(z).LengthSquareDirected(z)
	pr3 := t.P3.ToVector().ProjectOn(z).LengthSquareDirected(z)
	return math.Min(pr1, math.Min(pr2, pr3)), math.Max(pr1, math.Max(pr2, pr3))
}
