package slicers

import (
	"github.com/l1va/goosli"
	"github.com/l1va/goosli/commands"
	"math"
)

func SliceWithSlope(mesh goosli.Mesh) []commands.Layer {

	bb := mesh.BoundingBox()

	c := goosli.Point{X: (bb.MaxX - bb.MinX) / 2, Y: (bb.MaxY - bb.MinY) / 2, Z: bb.MinZ} // let it be origin
	//c := goosli.Point{X: 0, Y: 0, Z: 0} // let it be origin

	alpha := 30 * math.Pi / 180
	// transposed matrix to rotate around X
	mx := goosli.V(1, 0, 0)
	my := goosli.V(0, math.Cos(alpha), math.Sin(alpha))
	mz := goosli.V(0, -math.Sin(alpha), math.Cos(alpha))

	triangles := make([]goosli.Triangle, len(mesh.Triangles))
	rotatedMesh := goosli.NewMesh(triangles)
	for i, t := range mesh.Triangles {
		p1 := c.VectorTo(t.P1).Rotate(mx, my, mz).ToPoint(c)
		p2 := c.VectorTo(t.P2).Rotate(mx, my, mz).ToPoint(c)
		p3 := c.VectorTo(t.P3).Rotate(mx, my, mz).ToPoint(c)
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}

	cmds := Slice3DOF(rotatedMesh)

	// Reverse rotation
	ralpha := -30 * math.Pi / 180
	// transposed matrix to rotate around X
	rmx := goosli.V(1, 0, 0)
	rmy := goosli.V(0, math.Cos(ralpha), math.Sin(ralpha))
	rmz := goosli.V(0, -math.Sin(ralpha), math.Cos(ralpha))

	res := make([]commands.Layer, len(cmds))
	for i, cm := range cmds {
		paths := make([]commands.Path, len(cm.Paths))
		for j, p := range cm.Paths {
			lines := make([]commands.Line, len(p.Lines))
			for k, line := range p.Lines {
				lines[k].P1 = c.VectorTo(line.P1).Rotate(rmx, rmy, rmz).ToPoint(c)
				lines[k].P2 = c.VectorTo(line.P2).Rotate(rmx, rmy, rmz).ToPoint(c)
			}
			paths[j].Lines = lines
		}
		res[i].Paths = paths
	}
	return res

}
