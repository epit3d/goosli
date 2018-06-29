package primitives

import (
	"math"
)

var (
	OriginPoint = Point{0, 0, 0}
)

type Point struct {
	X, Y, Z float64
}

func (a Point) Equal(b Point) bool {

	return math.Abs(a.X-b.X) < AlmostZero && math.Abs(a.Y-b.Y) < AlmostZero && math.Abs(a.Z-b.Z) < AlmostZero
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

func (a Point) DistanceToLine(b Point, c Point) float64 {
	return a.ProjectOnLine(b, c).DistanceTo(a)
}

func (a Point) ProjectOnLine(b, c Point) Point {
	ba := b.VectorTo(a)
	bcUnit := b.VectorTo(c).Normalize()
	return b.Shift(bcUnit.MulScalar(ba.Dot(bcUnit)))
}

// Crossing Number method: calculates count of intersections plane with polygon(path) in one direction,
// if odd - inside. Can have problems with difficult forms as shown here
// http://geomalgorithms.com/a03-_inclusion.html
func (p Point) Inside(path Path) bool {
	if len(path.Lines) == 0 {
		return false
	}
	n := path.Lines[0].P1.VectorTo(path.Lines[0].P2)
	pl := Plane{p, n}
	v := p.VectorTo(path.Lines[0].P1).ProjectOnPlane(pl)
	c := 0
	for i := 0; i < len(path.Lines); i++ {
		line := path.Lines[i]
		inters := pl.IntersectSegment(line.P1, line.P2)
		if inters != nil && p.VectorTo(*inters).CodirectedWith(v) {
			c += 1
		}
	}
	return c%2 != 0
}
