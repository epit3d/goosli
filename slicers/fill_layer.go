package slicers

import (
	"math"
	"sort"

	. "github.com/l1va/goosli/primitives"
)

func calcFillPlanesCommon(mesh *Mesh, settings Settings, MyAxis Vector, step float64) []Plane {
	minx, maxx := mesh.MinMaxZ(MyAxis)
	//get minx which is  "k * step"
	minx_step := float64(math.Floor((minx)/step)) * step
	n := int(math.Ceil((maxx - minx_step) / step))
	curP := MyAxis.MulScalar(minx_step).ToPoint()

	planes := []Plane{}
	for i := 0; i < n; i++ {
		curP = curP.Shift(MyAxis.MulScalar(step))
		planes = append(planes, Plane{curP, MyAxis})
	}

	return planes
}

// filling with lines
func CalcFillPlanesLines(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity)) * settings.LineWidth
	planes := calcFillPlanesCommon(mesh, settings, AxisX, step)

	return planes
}

// filling with squares
func CalcFillPlanesSquares(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity/2)) * settings.LineWidth

	//planes := calcFillPlanesCommon(mesh, settings, V(1, 0, 0), step)
	//planes = append(planes, CalcFillPlanes0(mesh, settings, V(0, 1, 0), step)...)

	planes := calcFillPlanesCommon(mesh, settings, V(0.7071067812, 0.7071067812, 0), step)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(-0.7071067812, 0.7071067812, 0), step)...)

	return planes
}

// filling with triangles
func CalcFillPlanesTriangles(mesh *Mesh, settings Settings) []Plane {
	step := (100 / float64(settings.FillDensity/3)) * settings.LineWidth

	planes := calcFillPlanesCommon(mesh, settings, V(1, 0, 0), step)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(0.5, 0.8660254038, 0), step)...)
	planes = append(planes, calcFillPlanesCommon(mesh, settings, V(-0.5, 0.8660254038, 0), step)...)

	return planes
}

func CalcFillPlanes(mesh *Mesh, settings Settings) []Plane {
	planes := []Plane{}
	switch settings.FillingType {
	case "Lines":
		planes = CalcFillPlanesLines(mesh, settings)
	case "Squares":
		planes = CalcFillPlanesSquares(mesh, settings)
	case "Triangles":
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

func intersectByPlanePolygonWise(pathes []Path, plane Plane) []Path {
	// TODO: Merge polygons (some polygons may be inside other ones)
	outPaths := []Path{}
	for _, pth := range pathes {
		// Skip polylines
		if !pth.IsClosed() {
			continue
		}
		// Skip holes
		if pth.IsHole() {
			continue
		}

		// Find intersection points with every segment
		pts   := []Point{}
		for i := 1; i < len(pth.Points); i++ {
			p := plane.IntersectSegment(pth.Points[i-1], pth.Points[i])
			if p != nil {
				pts = append(pts, *p)
			}
		}
		if len(pts) >= 2 {
			// Sort intersection points by one of the axes
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

			// Add pairs of points
			for i := 1; i < len(pts); i += 2 {
				outPaths = append(outPaths, Path{Points: []Point{pts[i-1], pts[i]}})
			}
		}
	}

	// Return filling paths if found
	if len(outPaths) > 0 {
		return outPaths
	} else {
		return nil
	}
}

func intersectByPlane(pathes []Path, plane Plane) []Path {
	pts := []Point{}
	for _, pth := range pathes {
		for i := 1; i < len(pth.Points); i++ {

			p := plane.IntersectSegment(pth.Points[i-1], pth.Points[i])
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

	if len(pts) > 4 {
		paths := []Path{}
		for i := 1; i < len(pts); i += 2 {
			paths = append(paths, Path{Points: []Point{pts[i-1], pts[i]}})
		}
		return paths
	}
	if len(pts) == 2 || len(pts) == 3 {
		return []Path{{Points: []Point{pts[0], pts[len(pts)-1]}}}
	}
	if len(pts) == 4 {
		return []Path{{Points: []Point{pts[0], pts[1]}}, {Points: []Point{pts[2], pts[3]}}}
	}
	return nil //can not happen
}

func nearAngle(val, near float64) bool {
	return val < near+3 && val > near-3
}
