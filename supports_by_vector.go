package helpers

import (
	"log"
	"math"
	. "github.com/l1va/goosli/primitives"
)


func MakeSupportLines(t Triangle, a Plane) []Line {

	p1 := ProectionPointToPlane(t.P1, a)
	p2 := ProectionPointToPlane(t.P2, a)
	p3 := ProectionPointToPlane(t.P3, a)

	var	lines []Line;
	lines = append(lines, Line{P1: p1, P2: p2})
	lines = append(lines, Line{P1: p1, P2: p3})
	lines = append(lines, Line{P1: p2, P2: p3})

	return lines

}

func ProectionTriangleToPlane(t *Triangle, a Plane) Triangle {

	p1 := ProectionPointToPlane(t.P1, a)
	p2 := ProectionPointToPlane(t.P2, a)
	p3 := ProectionPointToPlane(t.P3, a)

	return Triangle{N: a.N, P1: p1, P2: p2, P3:p3}
}

func ProectionPointToPlane(M Point, a Plane)  Point {

	var x,y,z float64

	N := a.N
	P := a.P
	
	//plane equation
	//N.X*(x - P.X) + N.Y*(y - P.Y) + N.Z*(z - P.Z) = 0
	// canonical straight line equations
	//(x - M.X)/N.X = (y - M.Y)/N.Y = (z - M.Z)/N.Z = L

	// optimized equations
	// 	x = L*N.X + M.X
	//  y = L*N.Y + M.Y
	//  z = L*N.Z + M.Z
	//  N.X*x + N.Y*y + N.Z*z = N.X*P.X + N.Y*P.Y + N.Z*P.Z

	//solution by hand
	L := (N.X*(P.X - M.X) + N.Y*(P.Y - M.Y) + N.Z*(P.Z - M.Z))/(N.X*N.X + N.Y*N.Y + N.Z*N.Z)
	if math.IsNaN(L) {
		log.Fatal("Lambda = Nan")
	}

	x = L*N.X + M.X
	y = L*N.Y + M.Y
	z = L*N.Z + M.Z
	
	return Point{X: x, Y: y, Z: z}

}

