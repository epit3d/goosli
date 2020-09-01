package vip

import (
	//. "github.com/l1va/goosli/slicers"
	"math"

	. "github.com/l1va/goosli/primitives"
)

func MakeOffset(pth Path, nozzle float64, norm Vector) Path {
	if len(pth.Points) < 3 {
		println("can't make offset, empty path: n = ",  len(pth.Points))
		return pth //TODO: what return ?
	}
	norm = norm.Reverse()
	//norm := pth.Lines[0].ToVector().Cross(pth.Lines[len(pth.Lines)/2].ToVector()) //TODO: choose normal or have in path

	newPth := []Line{}
	for i := 1; i < len(pth.Points); i++ {
		p1 := pth.Points[i-1]
		p2 := pth.Points[i]
		bef := norm.Cross(ToVector(p1,p2))
		offv := bef.Normalize().MulScalar(nozzle)
		newPth = append(newPth, Line{p1.Shift(offv),p2.Shift(offv)})
	}

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
	}
	if closedPath {
		newPth[0] = newPth[len(newPth)-1]
		newPth = newPth[:len(newPth)-1] //remove added line
	}

	res := Path{}
	res.Points = append(res.Points, newPth[0].P1)
	for i:=0;i<len(newPth);i++{
		res.Points = append(res.Points, newPth[i].P2)
	}

	return res
}