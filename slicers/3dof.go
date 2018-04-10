package slicers

import (
	"github.com/l1va/goosli"
	"github.com/l1va/goosli/commands"
	"sort"
	"math"
)

func Slice3DOF(mesh goosli.Mesh) []commands.Layer {

	thickness := 0.2 //TODO: move to appropriate place

	bb := mesh.BoundingBox()

	triangles := mesh.CopyTriangles()

	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].MinZ < triangles[j].MinZ
	})

	n := int(math.Ceil((bb.MaxZ - bb.MinZ) / thickness))

	in := make(chan job, n)
	out := make(chan commands.Layer, n)

	goosli.DoInParallel(slicingWorker(in, out))

	index := 0
	var active []*goosli.Triangle
	for i := 0; i < n; i++ {
		z := goosli.RoundPlaces(bb.MinZ+thickness*float64(i)+thickness/2, 8)
		// remove triangles below plane
		newActive := active[:0]
		for _, t := range active {
			if t.MaxZ >= z {
				newActive = append(newActive, t)
			}
		}
		active = newActive
		// add triangles above plane
		for index < len(triangles) && triangles[index].MinZ <= z {
			active = append(active, triangles[index])
			index++
		}
		// copy triangles for worker job
		activeCopy := make([]*goosli.Triangle, len(active))
		copy(activeCopy, active)
		in <- job{z, activeCopy}
	}
	close(in)

	// read results from workers
	layers := make([]commands.Layer, n)
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
	z         float64
	triangles []*goosli.Triangle
}

func slicingWorker(in chan job, out chan commands.Layer) func(wi, wn int) {
	return func(_, _ int) {
		var paths []commands.Path
		for job := range in {
			paths = paths[:0]
			for _, t := range job.triangles {
				if line := intersectTriangle(job.z, t); line != nil {
					paths = append(paths, commands.Path{Lines: []commands.Line{*line}})
				}
			}
			out <- commands.Layer{Order: job.z, Paths: commands.JoinPaths(paths)}
		}
	}
}

func intersectSegment(z float64, p1, p2 goosli.Point) *goosli.Point {
	if p1.Z == p2.Z {
		return nil
	}
	t := (z - p1.Z) / (p2.Z - p1.Z)
	if t < 0 || t > 1 {
		return nil
	}
	p := p1.Shift(p1.VectorTo(p2).MulScalar(t))
	return &p
}

func intersectTriangle(z float64, t *goosli.Triangle) *commands.Line {
	v1 := intersectSegment(z, t.P1, t.P2)
	v2 := intersectSegment(z, t.P2, t.P3)
	v3 := intersectSegment(z, t.P3, t.P1)
	var p1, p2 goosli.Point
	if v1 != nil && v2 != nil {
		p1, p2 = *v1, *v2
	} else if v1 != nil && v3 != nil {
		p1, p2 = *v1, *v3
	} else if v2 != nil && v3 != nil {
		p1, p2 = *v2, *v3
	} else {
		return nil
	}
	p1 = p1.RoundPlaces(8)
	p2 = p2.RoundPlaces(8)
	if p1 == p2 {
		return nil

	}
	n := goosli.V(p1.Y-p2.Y, p2.X-p1.X, 0)
	if n.Dot(t.N) < 0 { // orientation according to triangle plane
		return &commands.Line{p1, p2}
	}
	return &commands.Line{p2, p1}
}
