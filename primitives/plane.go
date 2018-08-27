package primitives

var (
	PlaneXY = Plane{OriginPoint, AxisZ}
	PlaneYZ = Plane{OriginPoint, AxisX}
	PlaneXZ = Plane{OriginPoint, AxisY}
)

type Plane struct {
	P Point
	N Vector
}

func (p Plane) IntersectMesh(mesh *Mesh) Layer {
	var paths []Path
	for _, t := range mesh.Triangles {
		if line := p.IntersectTriangle(&t); line != nil {
			paths = append(paths, Path{Lines: []Line{*line}})
		}
	}
	return Layer{Order: 0, Paths: JoinPaths(paths)}
}

func (p Plane) Intersect(t *Triangle) bool {
	if t == nil {
		return false
	}
	min, max := t.MinMaxZ(p.N)
	pp := p.P.ToVector().Dot(p.N)
	return min <= pp && pp <= max
}

func (p Plane) IntersectTriangle(t *Triangle) *Line {
	if t == nil {
		return nil
	}

	v1 := p.IntersectSegment(t.P1, t.P2)
	v2 := p.IntersectSegment(t.P2, t.P3)
	v3 := p.IntersectSegment(t.P3, t.P1)
	var p1, p2 Point

	if v1 != nil && v2 != nil && !v1.Equal(*v2) {
		p1, p2 = *v1, *v2
	} else if v1 != nil && v3 != nil && !v1.Equal(*v3) {
		p1, p2 = *v1, *v3
	} else if v2 != nil && v3 != nil && !v2.Equal(*v3) {
		p1, p2 = *v2, *v3
	} else {
		return nil
	}
	p1 = p1.RoundPlaces(8)
	p2 = p2.RoundPlaces(8)
	if p1.Equal(p2) {
		return nil
	}
	n := p1.VectorTo(p2).Cross(p.N)
	if n.Dot(t.N) < 0 { // orientation according to triangle plane (against t normal)
		return &Line{p1, p2}
	}
	return &Line{p2, p1}
}

func (p Plane) IntersectSegment(p1, p2 Point) *Point {
	pr1 := p1.ToVector().Dot(p.N)
	pr2 := p2.ToVector().Dot(p.N)
	if pr1 == pr2 {
		return nil
	}
	z := p.P.ToVector().Dot(p.N)
	t := (z - pr1) / (pr2 - pr1)
	if t < 0 || t > 1 {
		return nil
	}
	res := p1.Shift(p1.VectorTo(p2).MulScalar(t))
	return &res
}

func (p Plane) IntersectPathCodirectedWith(path Path, v Vector) *Point {
	cp := FindCentroid(path)
	for _, line := range path.Lines { //TODO: optimize me
		pi := p.IntersectSegment(line.P1, line.P2)
		if pi != nil && cp.VectorTo(*pi).CodirectedWith(v) {
			return pi
		}
	}
	return nil
}

func (p Plane) PointInFront(v Point) bool {
	return p.P.VectorTo(v).Dot(p.N) >= 0
}
