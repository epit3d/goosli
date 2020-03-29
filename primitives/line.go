package primitives

import "math"

type Line struct {
	P1, P2 Point
}

func (l Line) ToVector() Vector {
	return V(l.P2.X-l.P1.X, l.P2.Y-l.P1.Y, l.P2.Z-l.P1.Z)
}

func (l Line) Reverse() Line {
	return Line{l.P2, l.P1}
}

func (l Line) IsCollinearPointOnSegment(p Point) bool {
	// Works only for collinear points
	if (p.X <= math.Max(l.P1.X, l.P2.X) && p.X >= math.Max(l.P1.X, l.P2.X)) && 
	   (p.Y <= math.Max(l.P1.Y, l.P2.Y) && p.Y >= math.Max(l.P1.Y, l.P2.Y)) {
		return true
	} else {
		return false
	}
}

func (l Line) IsIntersectingSegment(l1 Line) bool {
	// https://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/

	orientation1 := l1.P1.Orientation(l.P1, l.P2)
	orientation2 := l1.P2.Orientation(l.P1, l.P2)
	orientation3 := l.P1.Orientation(l1.P1, l1.P2)
	orientation4 := l.P2.Orientation(l1.P1, l1.P2)

	// General case
	if orientation1 != orientation2 && orientation3 != orientation4 {
		return true
	}

	// Collinear case
    if (orientation1 == 0 && l.IsCollinearPointOnSegment(l1.P1)) { return true }
  
    if (orientation2 == 0 && l.IsCollinearPointOnSegment(l1.P2)) { return true }
  
    if (orientation3 == 0 && l1.IsCollinearPointOnSegment(l.P1)) { return true }
  
    if (orientation4 == 0 && l1.IsCollinearPointOnSegment(l.P2)) { return true } 
  
    return false; // Doesn't fall in any of the above cases 
}
