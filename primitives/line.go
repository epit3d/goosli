package primitives

type Line struct {
	P1, P2 Point
}

func (l Line) ToVector() Vector {
	return V(l.P2.X-l.P1.X, l.P2.Y-l.P1.Y, l.P2.Z-l.P1.Z)
}

func (l Line) Reverse() Line {
	return Line{l.P2, l.P1}
}

func (l Line) IntersectLine(l2 *Line) *Point {

	a := l.ToVector()
	b := l2.ToVector()
	den := a.X + a.Y
	if den != 0 {
		m := (l2.P1.X + l2.P1.Z - l.P1.X - l.P1.Z - (a.X + a.Z)*(l2.P1.X + l2.P1.Y - l.P1.X - l.P1.Y)/den)
		m = m/((a.X + a.Z)*(b.X + b.Y)/den - b.X - b.Z)

		x := l2.P1.X + b.X*m
		y := l2.P1.Y + b.Y*m
		z := l2.P1.Z + b.Z*m

		return &Point{X: x, Y: y, Z:z}
	}
	den = a.X + a.Z
	if den != 0 {
		m := l2.P1.X + l2.P1.Y - l.P1.X - l.P1.Y - (a.X + a.Y)*(l2.P1.X + l2.P1.Z - l.P1.X - l.P1.Z)/den
		m = m/((a.X + a.Y)*(b.X + b.Z)/den - b.X - b.Y)

		x := l2.P1.X + b.X*m
		y := l2.P1.Y + b.Y*m
		z := l2.P1.Z + b.Z*m

		return &Point{X: x, Y: y, Z:z}
	}

	return nil
}
