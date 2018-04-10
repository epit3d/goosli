package goosli

import (
	"math"
)

type Mesh struct {
	Triangles []Triangle
}

func NewMesh(triangles []Triangle) Mesh {
	return Mesh{Triangles: triangles}
}

func (m *Mesh) MinZ() float64 {
	minz := math.MaxFloat64
	for _, t := range m.Triangles {
		minz = math.Min(minz, t.MinZ)
	}
	return minz
}

func (m *Mesh) MaxZ() float64 {
	maxz := -math.MaxFloat64
	for _, t := range m.Triangles {
		maxz = math.Max(maxz, t.MaxZ)
	}
	return maxz
}

func (m *Mesh) CopyTriangles() []*Triangle {
	triangles := make([]*Triangle, len(m.Triangles))

	work :=  func(wi, wn int) {
		for i := wi; i < len(m.Triangles); i += wn {
			t := m.Triangles[i]
			triangles[i] = &t
		}
	}
	DoInParallelAndWait(work)

	return triangles
}
