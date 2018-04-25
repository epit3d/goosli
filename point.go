package goosli

import "math"

type Point struct {
	X, Y, Z float64
}

func (a Point) VectorTo(b Point) Vector {
	return V(b.X-a.X, b.Y-a.Y, b.Z-a.Z)
}

func (a Point) Shift(b Vector) Point {
	return Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Point) DistanceTo(b Point) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2) + math.Pow(a.Z-b.Z, 2))
}

func (a Point) RoundPlaces(n int) Point {
	x := RoundPlaces(a.X, n)
	y := RoundPlaces(a.Y, n)
	z := RoundPlaces(a.Z, n)
	return Point{x, y, z}
}

func (a Point) ToVector() Vector {
	return V(a.X, a.Y, a.Z)
}

func (a Point) ToString() string {
	return "X" + StrF(a.X) + " Y" + StrF(a.Y) + " Z" + StrF(a.Z)
}
