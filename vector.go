package goosli

import "math"

type Vector struct {
	X, Y, Z float64
}

func V(x, y, z float64) Vector {
	return Vector{X: x, Y: y, Z: z}
}

func (a Vector) Cross(b Vector) Vector {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return V(x, y, z)
}

func (a Vector) Add(b Vector) Vector {
	return V(a.X+b.X, a.Y+b.Y, a.Z+b.Z)
}

func (a Vector) Sub(b Vector) Vector {
	return V(a.X-b.X, a.Y-b.Y, a.Z-b.Z)
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) MulScalar(b float64) Vector {
	return Vector{a.X * b, a.Y * b, a.Z * b}
}

func (a Vector) Rotate(mx, my, mz Vector) Vector {
	return V(a.Dot(mx), a.Dot(my), a.Dot(mz))
}

func (a Vector) ToPoint() Point {
	return Point{X: a.X, Y: a.Y, Z: a.Z}
}

func (a Vector) CodirectedWith(b Vector) bool {
	return a.Dot(b) >= 0
}

func (a Vector) Length() float64 {
	return math.Sqrt(a.LengthSquare())
}

func (a Vector) LengthSquare() float64 {
	return a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

func (a Vector) Normalize() Vector {
	n := a.Length()
	if n == 0 {
		return a
	}
	return V(a.X/n, a.Y/n, a.Z/n)
}
/*
func (a Vector) ProjectOn(b Vector) Vector {
	n := b.LengthSquare()
	if n == 0 {
		return b
	}
	return b.MulScalar(b.Dot(a) / n)
}*/
