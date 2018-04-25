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

func (m *Mesh) CopyTriangles() []*Triangle {
	triangles := make([]*Triangle, len(m.Triangles))

	work := func(wi, wn int) {
		for i := wi; i < len(m.Triangles); i += wn {
			t := m.Triangles[i]
			triangles[i] = &t
		}
	}
	DoInParallelAndWait(work)

	return triangles
}

func (m *Mesh) MinMaxZ(z Vector) (float64, float64) {
	if len(m.Triangles) == 0 {
		return 0, 0
	}
	minz := math.MaxFloat64
	maxz := -math.MaxFloat64

	for _, t := range m.Triangles {
		tminz, tmaxz := t.MinMaxZ(z)
		minz = math.Min(minz, tminz)
		maxz = math.Max(maxz, tmaxz)
	}

	return minz, maxz
}
