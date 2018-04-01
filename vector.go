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
