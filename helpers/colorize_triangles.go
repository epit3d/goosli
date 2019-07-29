package helpers

import (
	. "github.com/l1va/goosli/primitives"
)

//ColorizeTriangles - Returns indexes of triangles that has normal correlated with Z in angle
func ColorizeTriangles(m Mesh, angle float64) []bool {
	var res []bool
	for _, t := range m.Triangles {
		if t.N.Angle(AxisZ.Reverse()) < angle {
			res = append(res, true)
		} else {
			res = append(res, false)
		}
	}
	return res
}
