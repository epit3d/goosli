package slicers

import (
	"github.com/l1va/goosli"
	"github.com/fogleman/fauxgl"
)


func Crop(mesh *goosli.Mesh, plane goosli.Plane) *goosli.Mesh {
	var triangles []goosli.Triangle
	for _, t := range mesh.Triangles {
		f1 := p.pointInFront(t.V1.Position)
		f2 := p.pointInFront(t.V2.Position)
		f3 := p.pointInFront(t.V3.Position)
		if f1 && f2 && f3 {
			triangles = append(triangles, t)
		} else if f1 || f2 || f3 {
			triangles = append(triangles, p.clipTriangle(t)...)
		}
	}
	return goosli.NewMesh(triangles)
}
