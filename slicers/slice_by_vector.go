package slicers

import (
	"math"
	"sort"

	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
)

// SliceByVectorToGcode - Slicing on layers by vector Z
func SliceByVectorToGcode(mesh *Mesh, Z Vector, settings Settings) gcode.Gcode {
	layers := SliceByVector(mesh, Z, settings)
	layers = FillLayers(layers, CalcFillPlanes(mesh, settings))

	gcd := gcode.NewGcode(*settings.GcodeSettings)

	gcd.AddLayers(layers)
	return gcd
}

// SliceByVector - Slicing on layers by vector Z
func SliceByVector(mesh *Mesh, Z Vector, settings Settings) []Layer {

	thickness := settings.LayerHeight
	angle := settings.ColorizedAngle
	vn := settings.UnitVector
	supports := settings.SupportsOn
	support_offset := settings.SupportOffset

	if mesh == nil || len(mesh.Triangles) == 0 {
		return nil
	}
	Z = Z.Normalize()

	triangles := mesh.CopyTriangles()
	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].MinZ(Z) < triangles[j].MinZ(Z)
	})

	minz, maxz := mesh.MinMaxZ(Z)
	n := int(math.Ceil((maxz - minz) / thickness))

	sh := Z.MulScalar(maxz - minz).MulScalar(1.0 / float64(n))

	in := make(chan job, n)
	out := make(chan Layer, n)
	DoInParallel(slicingWorker(in, out))

	//make support triangles
	if supports {
		temp := mesh.CopyTriangles()
		sort.Slice(temp, func(i, j int) bool {
			return temp[i].MinZ(Z) > temp[j].MinZ(Z)
		})
		curPoint := Z.MulScalar(minz).ToPoint()
		plane := Plane{P: curPoint, N: Z}

		col_triangles := helpers.FilterTrianglesByColor(*mesh, angle, vn)

		if support_offset != 0.0 {
			var col_triangles_shifted []Triangle
			for _, t := range col_triangles {
				col_triangles_shifted = append(col_triangles_shifted, t.Shift(vn.MulScalar(-support_offset)))
			}
			col_triangles = col_triangles_shifted
		}

		col_lines := helpers.MakeUndoubledLinesFromTriangles(col_triangles)
		sup_lines := helpers.MakeSupportLines(col_lines, plane)
		var sup_triangles []Triangle
		for i := range col_lines {
			Tr1 := NewTriangle(sup_lines[i].P1, col_lines[i].P1, col_lines[i].P2)
			Tr2 := NewTriangle(sup_lines[i].P1, sup_lines[i].P2, col_lines[i].P2)

			lines := helpers.IntersectTriangles(Tr1, temp)
			if lines == nil {
				sup_triangles = append(sup_triangles, Tr1)
			}

			lines = helpers.IntersectTriangles(Tr2, temp)
			if lines == nil {
				sup_triangles = append(sup_triangles, Tr2)
			}
		}
		triangles = append(sup_triangles, triangles...)
	}

	index := 0
	var active []*Triangle
	curP := Z.MulScalar(minz).ToPoint().Shift(sh.MulScalar(0.5))
	for i := 0; i < n; i++ {
		plane := Plane{P: curP, N: Z}
		z := curP.ToVector().Dot(Z)
		// remove triangles below plane
		newActive := active[:0]
		for _, t := range active {
			if z <= t.MaxZ(Z) {
				newActive = append(newActive, t)
			}
		}
		active = newActive
		// add triangles above plane
		for index < len(triangles) && triangles[index].MinZ(Z) <= z {
			t := triangles[index]
			active = append(active, &t)
			index++
		}
		// copy triangles for worker job
		activeCopy := make([]*Triangle, len(active))
		copy(activeCopy, active)
		in <- job{order: i, plane: plane, triangles: activeCopy}
		curP = curP.Shift(sh)
	}
	close(in)

	// read results from workers
	layers := make([]Layer, n)
	for i := 0; i < n; i++ {
		layers[i] = <-out
	}

	// sort layers
	sort.Slice(layers, func(i, j int) bool {
		return layers[i].Order < layers[j].Order
	})

	// filter out empty layers
	if len(layers[0].Paths) == 0 {
		layers = layers[1:]
	}
	if len(layers[len(layers)-1].Paths) == 0 {
		layers = layers[:len(layers)-1]
	}

	return layers
}

type job struct {
	order     int
	plane     Plane
	triangles []*Triangle
}

func slicingWorker(in chan job, out chan Layer) func(wi, wn int) {
	return func(_, _ int) {
		var paths []Path
		for job := range in {
			paths = paths[:0]
			for _, t := range job.triangles {
				if line := job.plane.IntersectTriangle(t); line != nil {
					paths = append(paths, Path{Points: []Point{line.P1, line.P2}})
				}
			}
			out <- Layer{Order: job.order, Norm: job.plane.N, Paths: JoinPaths2(paths)}
		}
	}
}
