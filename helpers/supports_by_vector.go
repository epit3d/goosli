package helpers

import (
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

func IntersectTriangles(t Triangle, triangles []Triangle) []Line {

	var lines []Line
	for _, bt := range triangles {
		line := bt.IntersectTriangle(&t)
		if line != nil {
			lines = append(lines, *line)
		}
	}
	if len(lines) != 0 {
		return lines
	}
	return nil
}


func MakeSupportLines(lines []Line, a Plane) []Line{
	
	var pr_lines []Line
	for _,l := range lines {
		p1 := a.ProectionPointToPlane(l.P1)
		p2 := a.ProectionPointToPlane(l.P2)
		pr_lines = append(pr_lines, Line{P1: p1, P2: p2})
	}

	return pr_lines
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