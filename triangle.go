package goosli

type Triangle struct {
	N          Vector
	V1, V2, V3 Vector
}

func NewTriangle(v1, v2, v3 Vector) Triangle {
	n := normal(v1.Sub(v2), v2.Sub(v3))
	return Triangle{N: n, V1: v1, V2: v2, V3: v3}
}

func (t Triangle) fill(v1, v2, v3 Vector) {
	t.V1 = v1
	t.V2 = v2
	t.V3 = v3
	t.N = normal(v1.Sub(v2), v2.Sub(v3))
}

func normal(v1, v2 Vector) Vector {
	return v1.Cross(v2)
}
