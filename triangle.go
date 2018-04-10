package goosli

import (
	"math"
)

type Triangle struct {
	N          Vector
	P1, P2, P3 Point
	MinZ, MaxZ float64
}

func NewTriangle(p1, p2, p3 Point) Triangle {
	t := Triangle{}
	t.fill(p1,p2,p3)
	return t
}

func (t *Triangle) fill(p1, p2, p3 Point) {
	t.P1 = p1
	t.P2 = p2
	t.P3 = p3
	t.recalculate()
}

func (t *Triangle) recalculate(){
	t.N =  normal(t.P1.VectorTo(t.P2), t.P1.VectorTo(t.P3))
	t.MinZ = math.Min(math.Min(t.P1.Z, t.P2.Z), t.P3.Z)
	t.MaxZ = math.Max(math.Max(t.P1.Z, t.P2.Z), t.P3.Z)
}

func normal(v1, v2 Vector) Vector {
	return v1.Cross(v2)
}

