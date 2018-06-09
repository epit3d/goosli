package slicers

import (
	"github.com/l1va/goosli"
)

//SimplifyLine - DouglasPeucker algorithm of line simplification
func SimplifyLine(ps []goosli.Point, epsilon float64) []goosli.Point {
	// Find the point with the maximum distance
	dmax := 0.0
	index := 0
	end := len(ps)
	for i := 1; i < end-1; i++ {
		d := ps[i].DistanceToLine(ps[0], ps[end-1])
		if d > dmax {
			index = i
			dmax = d
		}
	}
	// If max distance is greater than epsilon, recursively simplify
	if dmax > epsilon {
		res1 := SimplifyLine(ps[:index+1], epsilon)
		res2 := SimplifyLine(ps[index:], epsilon)
		return append(res1[:len(res1)-1], res2...)
	} else {
		return []goosli.Point{ps[0], ps[end-1]}
	}
	return nil

}
