package primitives

import (
	"math"
)

var (
	AxisX = V(1, 0, 0)
	AxisY = V(0, 1, 0)
	AxisZ = V(0, 0, 1)
)

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

func (a Vector) Rotate(rm RotationMatrix) Vector {
	return V(a.Dot(rm.X), a.Dot(rm.Y), a.Dot(rm.Z))
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

func (a Vector) Reverse() Vector {
	return V(-a.X, -a.Y, -a.Z)
}

func (a Vector) Angle(b Vector) float64 {
	v := a.Dot(b) / a.Length() / b.Length()
	if v < -1 || v > 1 {
		return 0
	}
	return math.Acos(v) * 180 / math.Pi
}

func (a Vector) ProjectOn(b Vector) Vector {
	bl := b.LengthSquare() // because result = b.MulScalar(projectLen/|b|), where projectLen = b.Dot(a)/|b|
	if bl == 0 {
		return b
	}
	return b.MulScalar(b.Dot(a) / bl)
}

func (a Vector) ProjectOnPlane(p Plane) Vector {
	return a.Add(a.ProjectOn(p.N).Reverse()) // substracting orthogonal to the plane component from a
}
