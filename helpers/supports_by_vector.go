package helpers

import (
	"log"
	"math"
	. "github.com/l1va/goosli/primitives"
)

func MakeUndoubledLinesFromTriangles(col_triangles []Triangle) []Line {
	
	var lines []Line
	for _, t := range col_triangles {
		lines = append(lines, Line{P1: t.P1, P2:t.P2})
		lines = append(lines, Line{P1: t.P1, P2:t.P3})
		lines = append(lines, Line{P1: t.P3, P2:t.P2})
	}

	var new_lines []Line 
	var lines2 []Line
	for i, l1 := range lines{
		fl := false
		lines2 = []Line{}
		lines2 = append(lines2, lines[:i]...)
		lines2 = append(lines2, lines[i+1:]...)
		for _, l2 := range lines2{
			if l1 == l2 || l1 == l2.Reverse() {
				fl = true
				break
			}
		}
		if !fl{
			new_lines = append(new_lines, l1) 
		}
	}

	return new_lines
}

func MakeSupportLines(lines []Line, a Plane) []Line{
	
	var pr_lines []Line
	for _,l := range lines {
		p1 := ProectionPointToPlane(l.P1, a)
		p2 := ProectionPointToPlane(l.P2, a)
		pr_lines = append(pr_lines, Line{P1: p1, P2: p2})
	}

	return pr_lines
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

func compareAndDeleteLineFromSlice(lines []Line, temp_l Line) []Line {
	fl := true
	for j, l := range lines{
		if temp_l == l || temp_l == l.Reverse() {
			fl = false
			lines[j] = lines[len(lines)-1] // Copy last element to index j.
			lines = lines[:len(lines)-1]   // Truncate slice.
			break
		}
	}

	if fl {
		lines = append(lines, temp_l)
	}

	return lines
}

func DeleteInternalLines(t_arr []Triangle) []Line {

	var	lines []Line

	for _, t := range t_arr {
		lines = compareAndDeleteLineFromSlice(lines, Line{P1: t.P1, P2: t.P2})
		lines = compareAndDeleteLineFromSlice(lines, Line{P1: t.P1, P2: t.P3})
		lines = compareAndDeleteLineFromSlice(lines, Line{P1: t.P3, P2: t.P2})
	}
	return lines
}