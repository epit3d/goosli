package primitives

import "sort"

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

/*func (p Path) Enclose() Path {
	s := len(p.Lines)
	if !p.Lines[s-1].P2.Equal(p.Lines[0].P1) {
		p.Lines = append(p.Lines, Line{p.Lines[s-1].P2, p.Lines[0].P1})
	}
	return p
}*/

func toslice(m map[Point]Path) []Path {
	values := make([]Path, 0, len(m))

	for _, v := range m {
		values = append(values, v)
	}
	return values
}

type sortp struct {
	P   Point
	Ind int
}

func JoinPaths3(p []Path) []Path {
	n := len(p)
	w := make([]int, n)
	a := make([]sortp, n*2)
	for i := 0; i < n; i++ {
		w[i] = i
		a[2*i] = sortp{
			P:   p[i].Points[0],
			Ind: i,
		}
		a[2*i+1] = sortp{
			P:   p[i].Points[len(p[i].Points)-1],
			Ind: i,
		}
	}
	sort.Slice(a, func(i, j int) bool {
		b := a[i].P
		c := a[j].P
		return b.X < c.X || (AlmostZero(b.X-c.X) && b.Y < c.Y) || (AlmostZero(b.X-c.X) && AlmostZero(b.Y-c.Y) && b.Z < c.Z)
	})

	for i := 1; i < len(a); i++ {
		b := a[i-1]
		c := a[i]
		if b.P.Equal(c.P) {
			j1 := b.Ind
			for w[j1] != j1 {
				j1 = w[j1]
			}
			p1 := p[j1]
			j2 := c.Ind
			for w[j2] != j2 {
				j2 = w[j2]
			}
			if j1 == j2 {
				continue
			}
			p2 := p[j2]
			res := tryJoin(p1, p2)
			if res != nil {
				p[j1] = *res
				w[j2] = j1
			}
		}
	}
	res := []Path{}
	for i := 0; i < n; i++ {
		if w[i] == i {
			res = append(res, p[i])
		}
	}
	return res
}

func tryJoin(p1 Path, p2 Path) *Path {
	if p1.Points[len(p1.Points)-1].Equal(p2.Points[0]) {
		p1.Points = append(p1.Points, p2.Points[1:]...)
		return &p1
	}
	if p2.Points[len(p2.Points)-1].Equal(p1.Points[0]) {
		p2.Points = append(p2.Points, p1.Points[1:]...)
		return &p2
	}

	if p1.Points[0].Equal(p2.Points[0]) {
		p1.Points = append(p1.Reverse().Points, p2.Points[1:]...)
		return &p1
	}

	if p1.Points[len(p1.Points)-1].Equal(p2.Points[len(p2.Points)-1]) {
		p1.Points = append(p1.Points, p2.Reverse().Points[1:]...)

	}

	return nil
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
			if p, ok := lookup[cur.Points[0].MapKey()]; ok {
				yes = true
				delete(lookup, cur.Points[0].MapKey())
				p.Points = append(p.Points, cur.Points[1:]...)
				paths[i] = p
			} else {
				if _, ok := lookup[cur.Points[len(cur.Points)-1].MapKey()]; ok {
					paths[i] = cur.Reverse()
				} else {
					lookup[cur.Points[len(cur.Points)-1].MapKey()] = cur
					i++
				}
			}
		}
		paths = toslice(lookup)

	}
	return paths
}

/*
func JoinPaths(paths []Path) []Path {
	lookup := make(map[Point]Path, len(paths))
	for _, path := range paths {
		lookup[path.Lines[0].P1] = path
	}
	var result []Path
	for len(lookup) > 0 {
		var v Point
		for v = range lookup {
			break
		}
		var path Path
		for {
			if p, ok := lookup[v]; ok {
				path.Lines = append(path.Lines, p.Lines[0])
				delete(lookup, v)
				v = p.Lines[len(p.Lines)-1].P2
			} else {
				break
			}
		}
		result = append(result, path)
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			jp := tryJoin(result[i], result[j])
			if jp != nil {
				result[i] = *jp
				result[j] = result[len(result)-1]
				result = result[:len(result)-1]
				i = -1
				break
			}
		}
	}
	// clean from small lines
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i].Lines); j++ {
			if result[i].Lines[j].ToVector().Length() < 0.1 { //TODO: hardcode
				if j != 0 {
					result[i].Lines[j-1].P2 = result[i].Lines[j].P2
				} else if j < len(result[i].Lines)-1 {
					result[i].Lines[j+1].P1 = result[i].Lines[j].P1
				}
				result[i].Lines = append(result[i].Lines[:j], result[i].Lines[j+1:]...)
				j--
			}
		}
		if len(result[i].Lines) == 0 {
			result[i] = result[len(result)-1]
			result = result[:len(result)-1]
			i--
		}
	}
	return result
}*/

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
