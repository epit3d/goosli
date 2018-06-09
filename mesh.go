package goosli

import (
	"math"
	"bytes"
	"github.com/constabulary/gb/testdata/src/c"
)

type Mesh struct {
	Triangles []Triangle
}

func NewMesh(triangles []Triangle) Mesh {
	return Mesh{Triangles: triangles}
}


func (m *Mesh) RotateX(angle int, around Point) *Mesh {

	alpha := math.Pi * float64(angle) / 180.0
	// transposed matrix to rotate around X
	mx := V(1, 0, 0)
	my := V(0, math.Cos(alpha), math.Sin(alpha))
	mz := V(0, -math.Sin(alpha), math.Cos(alpha))

	return m.Rotate(mx,my,mz, around)
}


func (m *Mesh) RotateZ(angle int, around Point) *Mesh {

	alpha := math.Pi * float64(angle) / 180.0
	// transposed matrix to rotate around X
	mx := V(math.Cos(alpha), math.Sin(alpha), 0)
	my := V(-math.Sin(alpha), math.Cos(alpha), 0)
	mz := V(0, 0, 1)

	return m.Rotate(mx,my,mz, around)
}

func (m *Mesh) Rotate(mx,my,mz Vector, around Point) *Mesh {
	cv := around.ToVector()
	triangles := make([]Triangle, len(m.Triangles))
	rotatedMesh := NewMesh(triangles)
	for i, t := range m.Triangles {
		p1 := around.VectorTo(t.P1).Rotate(mx, my, mz).Add(cv).ToPoint()
		p2 := around.VectorTo(t.P2).Rotate(mx, my, mz).Add(cv).ToPoint()
		p3 := around.VectorTo(t.P3).Rotate(mx, my, mz).Add(cv).ToPoint()
		rotatedMesh.Triangles[i].Fill(p1, p2, p3)
	}
	return &rotatedMesh
}


func (m *Mesh) Shift(v Vector) {
	if m == nil {
		return
	}
	for i, t:= range(m.Triangles){
		m.Triangles[i] = t.Shift(v)
	}
}

func (m *Mesh) CopyTriangles() []Triangle {
	triangles := make([]Triangle, len(m.Triangles))

	work := func(wi, wn int) {
		for i := wi; i < len(m.Triangles); i += wn {
			triangles[i] = m.Triangles[i]
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
