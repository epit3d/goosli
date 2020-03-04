package slicers

import (
	"github.com/l1va/goosli/debug"
	. "github.com/l1va/goosli/primitives"
	"math"
	"sort"
)

func calcFillPlanesCommon(mesh *Mesh, settings Settings, MyAxis Vector, step float64) []Plane {
	minx, maxx := mesh.MinMaxZ(MyAxis)
	curP := Point{X: 0, Y: 0, Z: 0}
	// Todo do not make planes outside the body
	nMax := int(math.Ceil((maxx - 0) / step))
	nMin := int(math.Ceil((0 - minx) / step))
	planes := []Plane{}
	planes = append(planes, Plane{curP, MyAxis})

	for i := 0; i < nMax; i++ {
		curP = curP.Shift(MyAxis.MulScalar(step))
		planes = append(planes, Plane{curP, MyAxis})
	}

	curP = Point{X: 0, Y: 0, Z: 0}
	for i := 0; i < nMin; i++ {
		curP = curP.Shift(MyAxis.MulScalar(-step))
		planes = append(planes, Plane{curP, MyAxis})
	}
	return planes
}

// filling with lines
func CalcFillPlanesLines(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity)) * settings.Nozzle
	planes := calcFillPlanesCommon(mesh, settings, AxisX, step)

	return planes
}

// filling with squares
func CalcFillPlanesSquares(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity/2)) * settings.Nozzle

	//planes := CalcFillPlanes0(mesh, settings, V(1, 0, 0), step)
	//planes = append(planes, CalcFillPlanes0(mesh, settings, V(0, 1, 0), step)...)

	planes := calcFillPlanesCommon(mesh, settings, V(0.7071067812, 0.7071067812, 0), step)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(-0.7071067812, 0.7071067812, 0), step)...)

	return planes
}

// filling with triangles
func CalcFillPlanesTriangles(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity/3)) * settings.Nozzle

	planes := calcFillPlanesCommon(mesh, settings, V(1, 0, 0), step)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(0.5, 0.8660254038, 0), step)...)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(-0.5, 0.8660254038, 0), step)...)

	return planes
}

func CalcFillPlanes(mesh *Mesh, settings Settings) []Plane {
	planes := []Plane{}
	switch settings.FillingType {
	case 0:
		planes = CalcFillPlanesLines(mesh, settings)
	case 1:
		planes = CalcFillPlanesSquares(mesh, settings)
	case 2:
		planes = CalcFillPlanesTriangles(mesh, settings)
	default:
		//	for future changings
		planes = CalcFillPlanesLines(mesh, settings)
	}
	return planes
}

func FillLayers(layers []Layer, planes []Plane) []Layer { //TODO: can be paralleled
	for i, layer := range layers {
		inner := layer.InnerPs
		if inner == nil { //if one layer only
			inner = layer.Paths
			//println("inner nil")
		}
		for _, plane := range planes {
			pth := intersectByPlane(inner, plane)
			if pth != nil {
				layers[i].Fill = append(layers[i].Fill, pth...)
			}
		}
	}
	return layers
}

var (
	x = 0
)

func intersectByPlane(pathes []Path, plane Plane) []Path {
	pts := []Point{}
	for _, pth := range pathes {
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
		if x < 23 {
			for i := 1; i < len(pts); i += 2 {
				debug.AddLine(Line{pts[i-1], pts[i]}, debug.GreenColor)
			}
			x += 1
		}
		println("do not know how to fill, pts > 4, skipping :", len(pts))
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
