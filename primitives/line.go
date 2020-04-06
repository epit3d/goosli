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
	v := math.Sqrt(math.Pow(l.P2.X - l.P1.X, 2) +
				   math.Pow(l.P2.Y - l.P1.Y, 2) +
				   math.Pow(l.P2.Z - l.P1.Z, 2))
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

func (l Line) IntersectLine(l1 *Line) *Point {

	orientation1 := l1.P1.Orientation(l.P1, l.P2)
	orientation2 := l1.P2.Orientation(l.P1, l.P2)
	orientation3 := l.P1.Orientation(l1.P1, l1.P2)
	orientation4 := l.P2.Orientation(l1.P1, l1.P2)

	// General case
	if orientation1 != orientation2 && orientation3 != orientation4 {

		a := l.ToVector()
		b := l1.ToVector()

		var m  float64
		if a.X != 0 {
			m = (l1.P1.Y - l.P1.Y - (a.Y/a.X)*(l1.P1.X - l.P1.X)) / (b.X*a.Y/a.X - b.Y)
		} else if a.Y != 0 {
			m = (l1.P1.X - l.P1.X - (a.X/a.Y)*(l1.P1.Y - l.P1.Y)) / (a.X*b.Y/a.Y - b.X)
		} else if a.Z != 0 {
			m = (l1.P1.X - l.P1.X - (a.X/a.Z)*(l1.P1.Z - l.P1.Z)) / (a.X*b.Z/a.Z - b.X)
		}

		x := l1.P1.X + b.X*m
		y := l1.P1.Y + b.Y*m
		z := l1.P1.Z + b.Z*m

		return &Point{x, y, z}
	}

	// Collinear case
    if (orientation1 == 0 && l.IsCollinearPointOnSegment(l1.P1)) {
		println("l1.P1  ", l1.P1.Z)
		return &l1.P1 }

    if (orientation2 == 0 && l.IsCollinearPointOnSegment(l1.P2)) {
		println("l1.P2  ", l1.P2.Z)
		return &l1.P2 }

    if (orientation3 == 0 && l1.IsCollinearPointOnSegment(l.P1)) {
		println("l.P1  ", l.P1.Z)
		return &l.P1 }

    if (orientation4 == 0 && l1.IsCollinearPointOnSegment(l.P2)) {
		println("l2.P2  ", l.P2.Z)
		return &l.P2 }

    return nil; // Doesn't fall in any of the above cases

}
