package slicers

import (
	"sort"
	"math"
	. "github.com/l1va/goosli/primitives"
)

func CalcFillPlanes(mesh *Mesh, settings Settings) []Plane {
	minx, maxx := mesh.MinMaxZ(AxisX)
	step := (100 / float64(settings.FillDensity)) * settings.Nozzle
	curP := AxisX.MulScalar(minx).ToPoint()
	n := int(math.Ceil((maxx - minx) / step))

	planes := []Plane{}
	for i := 0; i < n; i++ {
		curP = curP.Shift(AxisX.MulScalar(step))
		planes = append(planes, Plane{curP, AxisX})
	}
	return planes
}

func FillLayers(layers []Layer, planes []Plane) []Layer {
	for i, layer := range layers {
		res := layer
		for _, plane := range planes {
			pth := intersectByPlane(layer, plane)
			if pth != nil {
				res.Paths = append(res.Paths, pth...)
			}
		}
		layers[i] = res
	}
	return layers
}

func intersectByPlane(layer Layer, plane Plane) []Path {
	pts := []Point{}
	for _, pth := range layer.Paths {
		for _, line := range pth.Lines {

			p := plane.IntersectSegment(line.P1, line.P2)
			if p != nil {
				pts = append(pts, *p)
			}
		}
	}
	if len(pts) < 2 {
		return nil
	}
	ang := plane.N.ProjectOnPlane(PlaneXY).Angle(AxisY)
	if nearAngle(ang, 0) || nearAngle(ang, 180) {
		sort.Slice(pts, func(i, j int) bool { //sort by X
			return pts[i].X < pts[j].X
		})
	} else {
		sort.Slice(pts, func(i, j int) bool { //sort by Y
			return pts[i].Y < pts[j].Y
		})
	}

	if len(pts) > 4 { //TODO: any ideas ?
		println("do not know how to fill, pts > 4, skipping, layer order: ", layer.Order)
		return nil
	}
	if len(pts) == 2 || len(pts) == 3 {
		return []Path{Path{Lines: []Line{Line{pts[0], pts[len(pts)-1]}}}}
	}
	if len(pts) == 4 {
		return []Path{Path{Lines: []Line{Line{pts[0], pts[1]}}}, Path{Lines: []Line{Line{pts[2], pts[3]}}}}
	}
	return nil //can not happen
}

func nearAngle(val, near float64) bool {
	return val < near+3 && val > near-3
}
