package vip

import (
	//. "github.com/l1va/goosli/slicers"
	. "github.com/l1va/goosli/primitives"
	"math"
)

func MakeOffset(pth Path, nozzle float64, norm Vector) Path {

	if len(pth.Points) < 3 {
		println("empty path: n = ",  len(pth.Points))
		return pth //TODO: what return ?
	}
	norm = norm.Reverse()
	//norm := pth.Lines[0].ToVector().Cross(pth.Lines[len(pth.Lines)/2].ToVector()) //TODO: choose normal or have in path

	newPth := []Line{}
	for i := 1; i < len(pth.Points); i++ {
		p1 := pth.Points[i-1]
		p2:=pth.Points[i]
		bef := norm.Cross(ToVector(p1,p2))
		offv := bef.Normalize().MulScalar(nozzle)
		newPth = append(newPth, Line{p1.Shift(offv),p2.Shift(offv)})
	}
	/*
		if level > 282 && level < 300 {
			println( pth.Lines[0].ToVector().Angle(pth.Lines[1].ToVector()))
			debug.AddLine(Line{pth.Lines[0].P1, pth.Lines[0].P1.Shift(norm.Normalize())})
			debug.AddLine(pth.Lines[0])
			debug.AddLine(pth.Lines[1])
			//debug.AddLine(Line{l1.P1, l1.P1.Shift(l1.ToVector().MulScalar(0.7))})
		}*/

	// find instersection of lines https://math.stackexchange.com/questions/270767/find-intersection-of-two-3d-lines

	closedPath := pth.Points[0].DistanceTo(pth.Points[len(pth.Points)-1]) < 0.001
	if closedPath {
		newPth = append(newPth, newPth[0]) // add first line to the end
	}
	for i := 0; i < len(newPth)-1; i++ {
		l1 := newPth[i]
		l2 := newPth[i+1]
		f := l2.ToVector()
		g := Line{l1.P1, l2.P2}.ToVector()
		e := l1.ToVector()
		a := f.Cross(g)
		b := f.Cross(e)
		c := 1.0
		if !a.CodirectedWith(b) {
			c = -1.0
		}
		inters := l1.P1.Shift(e.MulScalar(c * a.Length() / b.Length()))
		if math.IsNaN(inters.X) || math.IsNaN(inters.Y) || math.IsNaN(inters.Z) || inters.DistanceTo(l1.P2) > 0.2 || inters.DistanceTo(l2.P1) > 0.2 {
			inters = l1.P2
		}

		newPth[i].P2 = inters
		newPth[i+1].P1 = inters

		//if level == 282 {
		//	debug.AddLine(Line{l1.P1, l1.P1.Shift(e.MulScalar(0.5))})
		//}
		//if e.Length() < 0.4{
		//	count+=1
		//}
	}
	if closedPath {
		newPth[0] = newPth[len(newPth)-1]
		newPth = newPth[:len(newPth)-1] //remove added line
	}

	res := Path{}
	for i:=0;i<len(newPth);i++{
		res.Points = append(res.Points, newPth[i].P1)
	}
	res.Points = append(res.Points, newPth[len(newPth)-1].P2)

	return res
}

/*func MakeOffset2(pth Path, level int) Path {

	n := len(pth.Lines)
	if n < 3 {
		println("empty path: n = ", n)
		return pth //TODO: what return ?
	}

	if !pth.Lines[0].P1.Equal(pth.Lines[len(pth.Lines)-1].P2) {
		pth.Lines = append(pth.Lines, Line{pth.Lines[len(pth.Lines)-1].P2, pth.Lines[0].P1})
	}
	c := FindCentroid(pth)

	newPth := Path{}
	for i := 0; i < n; i++ {
		l1 := pth.Lines[i]
		coef := 0.05
		lx1 := Line{l1.P1, c}.ToVector().MulScalar(coef)
		lx2 := Line{l1.P2, c}.ToVector().MulScalar(coef)
		newl1 := Line{l1.P1.Shift(lx1), l1.P2.Shift(lx2)}
		newPth.Lines = append(newPth.Lines, newl1)

		if level == 285 || level == 286 {

			//debug.AddLine(Line{l1.P1, l1.P1.Shift(offv)})
		}
	}

	return newPth
}
*/