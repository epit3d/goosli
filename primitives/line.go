package primitives

import (
	"math"
)

type Line struct {
	P1, P2 Point
}

func (l Line) ToVector() Vector {
	return V(l.P2.X-l.P1.X, l.P2.Y-l.P1.Y, l.P2.Z-l.P1.Z)
}

func (l Line) Reverse() Line {
	return Line{l.P2, l.P1}
}

func (l Line) Len() float64 {
	v := math.Sqrt(math.Pow(l.P2.X-l.P1.X, 2) +
		math.Pow(l.P2.Y-l.P1.Y, 2) +
		math.Pow(l.P2.Z-l.P1.Z, 2))
	return v
}

func (l Line) IsCollinearPointOnSegment(p Point) bool {
	// Works only for collinear points
	if (p.X <= math.Max(l.P1.X, l.P2.X) && p.X >= math.Min(l.P1.X, l.P2.X)) &&
		(p.Y <= math.Max(l.P1.Y, l.P2.Y) && p.Y >= math.Min(l.P1.Y, l.P2.Y)) &&
		(p.Z <= math.Max(l.P1.Z, l.P2.Z) && p.Z >= math.Min(l.P1.Z, l.P2.Z)) {
		return true
	} else {
		return false
	}
}

func (l Line) IsIntersectingSegment(l1 *Line) bool {
	intersectionPoint := l.IntersectLine(l1)
	if intersectionPoint != nil {
		return true
	} else {
		return false
	}
}

func (l Line) IntersectLine(l1 *Line) *Point {

	orientation1 := l1.P1.Orientation(l.P1, l.P2)
	orientation2 := l1.P2.Orientation(l.P1, l.P2)
	orientation3 := l.P1.Orientation(l1.P1, l1.P2)
	orientation4 := l.P2.Orientation(l1.P1, l1.P2)

	// General case
	if orientation1 != orientation2 && orientation3 != orientation4 {

		a := l.ToVector()
		b := l1.ToVector()

		var m float64
		m = (l1.P1.Y - l.P1.Y - (a.Y/a.X)*(l1.P1.X-l.P1.X)) / (b.X*a.Y/a.X - b.Y)

		if math.IsNaN(m) {
			m = (l1.P1.X - l.P1.X - (a.X/a.Y)*(l1.P1.Y-l.P1.Y)) / (a.X*b.Y/a.Y - b.X)
		}
		if math.IsNaN(m) {
			m = (l1.P1.X - l.P1.X - (a.X/a.Z)*(l1.P1.Z-l.P1.Z)) / (a.X*b.Z/a.Z - b.X)
		}
		if math.IsNaN(m) {
			m = (l1.P1.Y - l.P1.Y - (a.Y/a.Z)*(l1.P1.Z-l.P1.Z)) / (a.Y*b.Z/a.Z - b.Y)
		}

		x := l1.P1.X + b.X*m
		y := l1.P1.Y + b.Y*m
		z := l1.P1.Z + b.Z*m

		return &Point{x, y, z}
	}

	// Collinear case
	if orientation1 == 0 && l.IsCollinearPointOnSegment(l1.P1) {
		println("l1.P1  ", l1.P1.Z)
		return &l1.P1
	}

	if orientation2 == 0 && l.IsCollinearPointOnSegment(l1.P2) {
		println("l1.P2  ", l1.P2.Z)
		return &l1.P2
	}

	if orientation3 == 0 && l1.IsCollinearPointOnSegment(l.P1) {
		println("l.P1  ", l.P1.Z)
		return &l.P1
	}

	if orientation4 == 0 && l1.IsCollinearPointOnSegment(l.P2) {
		println("l2.P2  ", l.P2.Z)
		return &l.P2
	}

	return nil // Doesn't fall in any of the above cases

}

func (l1 Line) IsCollinear(l2 Line) bool {
	v1 := l1.ToVector()
	v2 := l2.ToVector()
	res := v1.Cross(v2)
	return AlmostZero(res.Length())
}

// find instersection of lines https://math.stackexchange.com/questions/270767/find-intersection-of-two-3d-lines
func (l1 Line) IntersectLine2(l2 Line) *Point {
	f := l2.ToVector()
	g := Line{l1.P1, l2.P2}.ToVector()
	e := l1.ToVector()
	a := f.Cross(g)
	b := f.Cross(e)
	//println(a.Length(), b.Length())
	if AlmostZero(a.Length()) || AlmostZero(b.Length()) {
		println("A OR B ZERO", StrF(a.Length()), StrF(b.Length()), StrF(l1.Len()), StrF(l2.Len()))
		return nil
	}
	c := 1.0
	if !a.CodirectedWith(b) {
		c = -1.0
	}
	inters := l1.P1.Shift(e.MulScalar(c * a.Length() / b.Length()))
	if math.IsNaN(inters.X) || math.IsNaN(inters.Y) || math.IsNaN(inters.Z) || inters.DistanceTo(l1.P2) > 5 || inters.DistanceTo(l2.P1) > 5 { //TODO: distance?
		println("MORE 5", StrF(inters.DistanceTo(l1.P2)), StrF(inters.DistanceTo(l2.P1)))
		return nil
	}

	//mm := Line{l1.P1, inters}.ToVector()
	d1 := inters.DistanceToLine(l1.P1, l1.P2) //mm.Cross(e).Length() / e.Length()
	//mm = Line{l2.P1, inters}.ToVector()
	d2 := inters.DistanceToLine(l2.P1, l2.P2) //mm.Cross(f).Length() / f.Length()
	if AlmostZero(d1) && AlmostZero(d2) {
		return &inters
	}
	println("DIST", StrF(d1), StrF(d2))
	return nil
}
