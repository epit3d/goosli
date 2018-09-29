package helpers

import (
	. "github.com/l1va/goosli/primitives"
//	"fmt"
	"log"
)

//CutMesh mesh in two meshes, first that inFront of plane, second - outFront
func CutMesh(mesh *Mesh, p Plane) (*Mesh, *Mesh, error) {
	if mesh == nil || len(mesh.Triangles) == 0 {
		return &Mesh{}, &Mesh{}, nil //fmt.Errorf("mesh is empty, nothing to cut")
	}

	var up []Triangle
	var down []Triangle

	for _, t := range mesh.Triangles { //TODO: make it more beautiful
		var inFront []Point
		var outFront []Point
		if p.PointInFront(t.P1) {
			inFront = append(inFront, t.P1)
		} else {
			outFront = append(outFront, t.P1)
		}
		if p.PointInFront(t.P2) {
			inFront = append(inFront, t.P2)
		} else {
			outFront = append(outFront, t.P2)
		}
		if p.PointInFront(t.P3) {
			inFront = append(inFront, t.P3)
		} else {
			outFront = append(outFront, t.P3)
		}
		if len(inFront) == 3 {
			up = append(up, t)
		} else if len(inFront) == 2 {
			line := p.IntersectTriangle(&t)
			if line == nil {
				log.Printf("failed to intersect triangle by plane2, skip it: %v, %v\n", t, p)
				continue
			}

			ts := splitOnThree(outFront[0], *line, t)

			up = append(up, ts[1])
			up = append(up, ts[2])
			down = append(down, ts[0])

		} else if len(inFront) == 1 {
			line := p.IntersectTriangle(&t)
			if line == nil {
				log.Printf("failed to intersect triangle by plane1, skip it: %v, %v\n", t, p)
				continue
			}
			ts := splitOnThree(inFront[0], *line, t)

			up = append(up, ts[0])
			down = append(down, ts[1])
			down = append(down, ts[2])

		} else {
			down = append(down, t)
		}
	}

	//if len(up) == 0 || len(down) == 0 {
	//	return nil, nil, fmt.Errorf("one of meshes is empty")
	//}

	resUp := NewMesh(up)
	resDown := NewMesh(down)
	return &resUp, &resDown, nil
}

// p1 - point that from triangle 1
// ans consist of 3 triangle: 1 leaves as triangle, 2 and 3 - triangles from quadrangle
func splitOnThree(p1 Point, line Line, t Triangle) []Triangle {
	var ans []Triangle
	lp1 := line.P1
	lp2 := line.P2
	t1 := t.P1
	t2 := t.P2
	t3 := t.P3
	if t2.Equal(p1) {
		t1, t2, t3 = t2, t3, t1
	} else if t3.Equal(p1) {
		t1, t2, t3 = t3, t1, t2
	}

	if AlmostZero(t1.VectorTo(t2).Cross(t1.VectorTo(lp2)).Length() ){
		// lp2 lies on t1->t2 vector, but should lp1 lie
		lp1, lp2 = lp2, lp1
	}
	ans = append(ans, NewTriangle(t1, lp1, lp2))
	ans = append(ans, NewTriangle(lp1, t2, t3))
	ans = append(ans, NewTriangle(t3, lp2, lp1))
	return ans
}
