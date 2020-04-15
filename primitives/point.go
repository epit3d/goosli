package primitives

import (
	"fmt"
	"math"
)

var (
	OriginPoint = Point{0, 0, 0}
)

type Point struct {
	X, Y, Z float64
}

func (a Point) String() string {
	return fmt.Sprintf("X%s Y%s Z%s", StrF(a.X), StrF(a.Y), StrF(a.Z))
}
func (a Point) MapKey() Point {
	return a.RoundPlaces(8)
}

func (a Point) StringGcode(b Point) string {
	s_form := ""
	if math.Abs(a.X-b.X) >= 0.001 {
		s_form = fmt.Sprintf("%sX%s ", s_form, StrF(a.X))
	}

	if math.Abs(a.Y-b.Y) >= 0.001 {
		s_form = fmt.Sprintf("%sY%s ", s_form, StrF(a.Y))
	}

	if math.Abs(a.Z-b.Z) >= 0.001 {
		s_form = fmt.Sprintf("%sZ%s", s_form, StrF(a.Z))
	}
	return s_form
}

func (a Point) Equal(b Point) bool {

	return AlmostZero(a.X-b.X) && AlmostZero(a.Y-b.Y) && AlmostZero(a.Z-b.Z)
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
func (a Point) Inside(path Path) bool {
	if len(path.Points) < 4 {
		return false
	}
	n := path.Points[0].VectorTo(path.Points[1])
	pl := Plane{a, n}
	v := a.VectorTo(path.Points[0]).ProjectOnPlane(pl)
	c := 0
	for i := 1; i < len(path.Points); i++ {
		inters := pl.IntersectSegment(path.Points[i-1], path.Points[i])
		if inters != nil && a.VectorTo(*inters).CodirectedWith(v) {
			c += 1
		}
	}
	return c%2 != 0
}

// See https://www.geeksforgeeks.org/orientation-3-ordered-points/
func (r Point) Orientation(p Point, q Point) int {
    val := (q.Y - p.Y) * (r.X - q.X) - (q.X - p.X) * (r.Y - q.Y)

    if AlmostZero(val) {
		return 0 // Collinear
	} else if val > 0 {
		return 1 // Clockwise
	} else {
		return 2 // Counterclockwise
	}
}
