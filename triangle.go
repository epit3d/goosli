package goosli

import (
	"math"
)

type Triangle struct {
	N                                  Vector
	P1, P2, P3                         Point
	MinX, MaxX, MinY, MaxY, MinZ, MaxZ float64
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
	t.MinX = math.Min(math.Min(t.P1.X, t.P2.X), t.P3.X)
	t.MaxX = math.Max(math.Max(t.P1.X, t.P2.X), t.P3.X)
	t.MinY = math.Min(math.Min(t.P1.Y, t.P2.Y), t.P3.Y)
	t.MaxY = math.Max(math.Max(t.P1.Y, t.P2.Y), t.P3.Y)
	t.MinZ = math.Min(math.Min(t.P1.Z, t.P2.Z), t.P3.Z)
	t.MaxZ = math.Max(math.Max(t.P1.Z, t.P2.Z), t.P3.Z)
}

func normal(v1, v2 Vector) Vector {
	return v1.Cross(v2)
}
