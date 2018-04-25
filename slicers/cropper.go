package slicers

import (
	"github.com/l1va/goosli"
)

func Crop(mesh *goosli.Mesh, p goosli.Plane) *goosli.Mesh {
	var triangles []goosli.Triangle
	for _, t := range mesh.Triangles {
		var inFront []goosli.Point
		if p.PointInFront(t.P1) {
			inFront = append(inFront, t.P1)
		}
		if p.PointInFront(t.P2) {
			inFront = append(inFront, t.P2)
		}
		if p.PointInFront(t.P3) {
			inFront = append(inFront, t.P3)
		}
		if len(inFront) == 3 {
			triangles = append(triangles, t)
		} else if len(inFront) == 2 {
			line := p.IntersectTriangle(&t)
			triangles = append(triangles, goosli.NewTriangle(inFront[0], line.P1, line.P2))
			triangles = append(triangles, goosli.NewTriangle(inFront[0], line.P2, inFront[1]))
		} else if len(inFront) == 1 {
			line := p.IntersectTriangle(&t)
			triangles = append(triangles, goosli.NewTriangle(inFront[0], line.P1, line.P2))
		}
	}
	res := goosli.NewMesh(triangles)
	return &res
}
