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


type BoundingBox struct {
	MinX, MaxX, MinY,MaxY, MinZ, MaxZ float64
}

func (m *Mesh) BoundingBox() BoundingBox{
	minx := math.MaxFloat64
	maxx := -math.MaxFloat64
	miny := math.MaxFloat64
	maxy := -math.MaxFloat64
	minz := math.MaxFloat64
	maxz := -math.MaxFloat64

	for _, t := range m.Triangles {
		minx = math.Min(minx, t.MinX)
		maxx = math.Max(maxx, t.MaxX)
		miny = math.Min(miny, t.MinY)
		maxy = math.Max(maxy, t.MaxY)
		minz = math.Min(minz, t.MinZ)
		maxz = math.Max(maxz, t.MaxZ)
	}

	return BoundingBox{
		MinX: minx,
		MaxX: maxx,
		MinY: miny,
		MaxY: maxy,
		MinZ: minz,
		MaxZ: maxz,
	}
}