package slicers

import (
	"github.com/l1va/goosli"
	"fmt"
	"github.com/fogleman/simplify"
)

func Simplify(mesh *goosli.Mesh, resSize int) (*goosli.Mesh, error) {
	if mesh == nil {
		return nil, fmt.Errorf("cannot simplify, mesh is nil")
	}
	actualSize := len(mesh.Triangles)

	if actualSize <= resSize {
		return mesh, nil
	}

	factor := float64(resSize) / float64(actualSize)

	smesh := toFoglemanMesh(mesh)
	smesh = smesh.Simplify(factor)
	return toGoosliMesh(smesh), nil
}

func toFoglemanMesh(mesh *goosli.Mesh) *simplify.Mesh {
	var sts []*simplify.Triangle
	for _, t := range (mesh.Triangles) {
		v1 := simplify.Vector{t.P1.X, t.P1.Y, t.P1.Z}
		v2 := simplify.Vector{t.P2.X, t.P2.Y, t.P2.Z}
		v3 := simplify.Vector{t.P3.X, t.P3.Y, t.P3.Z}
		st := simplify.Triangle{v1, v2, v3}
		sts = append(sts, &st)
	}
	res := simplify.Mesh{sts}
	return &res
}
func toGoosliMesh(mesh *simplify.Mesh) *goosli.Mesh {
	var sts []goosli.Triangle
	for _, t := range (mesh.Triangles) {
		v1 := goosli.Point{t.V1.X, t.V1.Y, t.V1.Z}
		v2 := goosli.Point{t.V2.X, t.V2.Y, t.V2.Z}
		v3 := goosli.Point{t.V3.X, t.V3.Y, t.V3.Z}
		st := goosli.NewTriangle(v1, v2, v3)
		sts = append(sts, st)
	}
	res := goosli.Mesh{sts}
	return &res
}
