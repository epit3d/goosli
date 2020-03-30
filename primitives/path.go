package primitives

type Path struct {
	Points             []Point
	Retraction         bool
	RetractionSpeed    int
	RetractionDistance float64
}

func (p Path) Reverse() Path {
	res := Path{}
	for i := len(p.Points) - 1; i >= 0; i-- {
		res.Points = append(res.Points, p.Points[i])
	}
	return res
}

func (p Path) Equal(p2 Path) bool {
	if len(p.Points) != len(p2.Points) {
		return false
	}
	for i, pi := range p.Points {
		if !pi.Equal(p2.Points[i]) {
			return false
		}
	}
	return true
}

func (p Path) IsWindingCW() bool {
	// Assuming that the points lie in the same plane
	sum := 0.0
	for i := 1; i < len(p.Points); i++ {
		sum += (p.Points[i].X - p.Points[i-1].X) * (p.Points[i].Y + p.Points[i-1].Y)
	}

	sum += (p.Points[0].X - p.Points[len(p.Points) - 1].X) * (p.Points[0].Y + p.Points[len(p.Points) - 1].Y)

	if sum > 0.0 {
		return true
	} else if sum < 0.0 {
		return false
	} else {
		// Something went worng
		return false
	}
}

func (p Path) IsClosed() bool {
	if len(p.Points) <= 1 {
		return true
	}

	return p.Points[0].Equal(p.Points[len(p.Points) - 1])
}

func (p Path) IsHole() bool {
	return !p.IsWindingCW()
}

func (p Path) IsSolid() bool {
	return !p.IsHole()
}
 
func (p Path) Close() Path {
	res := Path{}
	for i := 0; i < len(p.Points); i++ {
		res.Points = append(res.Points, p.Points[i])
	}

	if !p.IsClosed() {
		res.Points = append(res.Points, p.Points[0])
	}

	return res
}

func (p Path) IsInside(p1 Path) bool {
	for i := 1; i < len(p.Points); i++ {
		for j := 1; j < len(p1.Points); j++ {
			l1 := Line{P1: p.Points[i-1], P2: p.Points[i]}
			l2 := Line{P1: p1.Points[j-1], P2: p1.Points[j]}

			if l1.IsIntersectingSegment(l2) {
				return false
			}
		}
	}

	for i := 0; i < len(p.Points); i++ {
		if p.Points[i].Inside(p1) {
			return true
		}
	}

	return false
}

func toslice(m map[Point]Path) []Path {
	values := make([]Path, 0, len(m))

	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func JoinPaths2(p []Path) []Path {
	paths := make([]Path, len(p))
	copy(paths, p)
	yes := true
	for yes {
		yes = false

		lookup := map[Point]Path{}
		for i := 0; i < len(paths); {
			cur := paths[i]
			if p, ok := lookup[cur.Points[0]]; ok {
				yes = true
				delete(lookup, cur.Points[0])
				p.Points = append(p.Points, cur.Points[1:]...)
				paths[i] = p
			} else {
				if _, ok := lookup[cur.Points[len(cur.Points)-1]]; ok {
					paths[i] = cur.Reverse()
				} else {
					lookup[cur.Points[len(cur.Points)-1]] = cur
					i++
				}
			}
		}
		paths = toslice(lookup)

	}
	return paths
}

func FindCentroid(path Path) Point { //TODO: refactorme
	a := 0.0
	a2 := 0.0
	a3 := 0.0
	cx := 0.0
	cy := 0.0
	cy2 := 0.0
	cz2 := 0.0
	cx2 := 0.0
	cz := 0.0
	for i := 1; i < len(path.Points); i++ {
		p := path.Points[i-1]
		p_1 := path.Points[i]
		d := p.X*p_1.Y - p_1.X*p.Y
		a += d
		cx += (p.X + p_1.X) * d
		cy += (p.Y + p_1.Y) * d

		d2 := p.X*p_1.Z - p_1.X*p.Z
		a2 += d2
		cx2 += (p.X + p_1.X) * d2
		cz += (p.Z + p_1.Z) * d2

		d3 := p.Y*p_1.Z - p_1.Y*p.Z
		a3 += d3
		cy2 += (p.Y + p_1.Y) * d3
		cz2 += (p.Z + p_1.Z) * d3
	}
	var x, x1, y, y1, z, z1 float64
	//println(a, a2, a3, "<>", cx, cx2, cy, cy2, cz, cz2, "\n")
	a = a * 0.5
	if !AlmostZero(a) {
		x = cx / 6 / a
		y = cy / 6 / a
	}
	a2 = a2 * 0.5
	if !AlmostZero(a2) {
		x1 = cx2 / 6 / a2
		z = cz / 6 / a2
	}
	a3 = a3 * 0.5
	if !AlmostZero(a3) {
		y1 = cy2 / 6 / a3
		z1 = cz2 / 6 / a3
	}

	//println("x:", x, " ", x1)
	//println("y:", y, " ", y1)
	//println("z:", z, " ", z1)

	if AlmostZero(a) && AlmostZero(a2) {
		x = path.Points[0].X
		y = y1
		z = z1
	}
	if AlmostZero(a) && AlmostZero(a3) {
		y = path.Points[0].Y
		x = x1
	}
	if AlmostZero(a2) && AlmostZero(a3) {
		z = path.Points[0].Z
	}

	if AlmostZero(a) && !AlmostZero(a2) {
		x = x1
	}
	if AlmostZero(a) && !AlmostZero(a3) {
		y = y1
	}
	if AlmostZero(a2) && !AlmostZero(a3) {
		z = z1
	}
	//println(x,x1,y,y1,z,z1,"\n")
	return Point{x, y, z}
}
