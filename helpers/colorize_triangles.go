package helpers

import (
	. "github.com/l1va/goosli/primitives"
)

//ColorizeTriangles - Returns true if normal of triangle correlated with Z less than in angle
func ColorizeTriangles(m Mesh, angle float64) []bool {
	var res []bool
	for _, t := range m.Triangles {
		if t.N.Angle(AxisZ.Reverse()) < 90 - angle {
			res = append(res, true)
		} else {
			res = append(res, false)
		}
	}
	return res
}
