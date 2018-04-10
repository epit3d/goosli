package goosli

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

func (a Vector) ToPoint(origin Point) Point {
	return Point{X: a.X + origin.X, Y: a.Y + origin.Y, Z: a.Z + origin.Z}
}
