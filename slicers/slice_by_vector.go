package slicers

import (
	. "github.com/l1va/goosli/primitives"
	"sort"
	"math"
	"bytes"
	"github.com/l1va/goosli/gcode"
)

// SliceByVectorToBuffer - Slicing on layers by vector Z
func SliceByVectorToBuffer(mesh *Mesh, Z Vector, settings Settings) bytes.Buffer {
	layers := SliceByVector(mesh, settings.LayerHeight, Z)
	settings.LayerCount = len(layers)
	smap := settings.ToMap()

	var buffer bytes.Buffer
	var cmds []gcode.Command
	cmds = append(cmds, gcode.LayersMoving{layers, 0})
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}

// SliceByVector - Slicing on layers by vector Z
func SliceByVector(mesh *Mesh, thickness float64, Z Vector) []Layer {

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
					paths = append(paths, Path{Lines: []Line{*line}})
				}
			}
			out <- Layer{Order: job.order, Paths: JoinPaths(paths)}
		}
	}
}
