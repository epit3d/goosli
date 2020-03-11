package helpers

import (
	. "github.com/l1va/goosli/primitives"
)

//ColorizeTriangles - Returns true if normal of triangle correlated with Z less than in angle
func ColorizeTriangles(m Mesh, angle float64, v Vector) []bool {
	var res []bool
	for _, t := range m.Triangles {
		if t.N.Angle(v.Reverse()) < 90-angle {
			res = append(res, true)
		} else {
			res = append(res, false)
		}
	}
	return res
}

func FilterTrianglesByColor(m Mesh, angle float64, v Vector) []Triangle {

	arr := ColorizeTriangles(m, angle, v)
	var colTriangles []Triangle
	for i, t := range m.Triangles {
		if arr[i] == true {
			colTriangles = append(colTriangles, t)
		}
	}
	return colTriangles
}