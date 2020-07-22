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

func (l Line) IsIntersectingSegment(l1 *Line) bool {
	intersectionPoint := l.IntersectLine(l1)
	if intersectionPoint != nil {
		return true
	} else {
		return false
	}
}

func (l Line) IntersectLine(l1 *Line) *Point {

	normal := make([]Vector, 4)
	orientation := make([]int, 4)
	var base_normal Vector

	normal[0] = l1.P1.VectorTo(l.P1).Cross(l1.P1.VectorTo(l.P2))
	normal[1] = l1.P2.VectorTo(l.P1).Cross(l1.P2.VectorTo(l.P2))
	normal[2] = l.P1.VectorTo(l1.P1).Cross(l.P1.VectorTo(l1.P2))
	normal[3] = l.P2.VectorTo(l1.P1).Cross(l.P2.VectorTo(l1.P2))

	for i, n := range normal {
		if AlmostZero(n.X)&&AlmostZero(n.Y)&&AlmostZero(n.Z) {
			orientation[i] = 0
			normal[i] = V(0,0,0)
		} else {
			base_normal = n
		}
	}

	for i, n := range normal {
		if base_normal.Dot(n) > 0 {
			orientation[i] = 1
		}
		if base_normal.Dot(n) < 0 {
			orientation[i] = 2
		}
	}

	// Collinear case
    if (orientation[0] == 0 && l.IsCollinearPointOnSegment(l1.P1)) {
		return &l1.P1 }

    if (orientation[1] == 0 && l.IsCollinearPointOnSegment(l1.P2)) {
		return &l1.P2 }

    if (orientation[2] == 0 && l1.IsCollinearPointOnSegment(l.P1)) {
		return &l.P1 }

    if (orientation[3] == 0 && l1.IsCollinearPointOnSegment(l.P2)) {
		return &l.P2 }

	// General case
	if orientation[0] != orientation[1] && orientation[2] != orientation[3] {

		a := l.ToVector()
		b := l1.ToVector()

		var m  float64
		m = (l1.P1.Y - l.P1.Y - (a.Y/a.X)*(l1.P1.X - l.P1.X)) / (b.X*a.Y/a.X - b.Y)

		if math.IsNaN(m) {
			m = (l1.P1.X - l.P1.X - (a.X/a.Y)*(l1.P1.Y - l.P1.Y)) / (a.X*b.Y/a.Y - b.X)
		}
		if math.IsNaN(m) {
			m = (l1.P1.X - l.P1.X - (a.X/a.Z)*(l1.P1.Z - l.P1.Z)) / (a.X*b.Z/a.Z - b.X)
		}
		if math.IsNaN(m) {
			m = (l1.P1.Y - l.P1.Y - (a.Y/a.Z)*(l1.P1.Z - l.P1.Z)) / (a.Y*b.Z/a.Z - b.Y)
		}

		x := l1.P1.X + b.X*m
		y := l1.P1.Y + b.Y*m
		z := l1.P1.Z + b.Z*m

		return &Point{x, y, z}
	}


    return nil; // Doesn't fall in any of the above cases

}

func (l Line) IntersectTriangle(t Triangle) bool {

	p1 := Line{P1: t.P1, P2: t.P2}.IntersectLine(&l)
	p2 := Line{P1: t.P1, P2: t.P3}.IntersectLine(&l)
	p3 := Line{P1: t.P2, P2: t.P3}.IntersectLine(&l)

	if p1 != nil || p2 != nil || p3 != nil {
		return true
	}

/*	if p1 == nil && p2 == nil && p3 == nil {
		if t.PointBelongs(l.P1) {  //or l.P2 - no difference
			return true
		}
	}*/

	return false

}