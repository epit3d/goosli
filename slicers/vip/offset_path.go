package vip

import (
	"github.com/l1va/goosli/debug"
	. "github.com/l1va/goosli/primitives"
)

func MakeOffset(pth Path, nozzle float64, norm Vector) Path {
	if len(pth.Points) < 2 {
		println("can't make offset, empty path: n = ", len(pth.Points))
		return pth //TODO: what return ?
	}
	norm = norm.Reverse()
	//norm := pth.Lines[0].ToVector().Cross(pth.Lines[len(pth.Lines)/2].ToVector()) //TODO: choose normal or have in path

	newPth := []Line{}
	for i := 1; i < len(pth.Points); i++ {
		p1 := pth.Points[i-1]
		p2 := pth.Points[i]
		bef := norm.Cross(ToVector(p1, p2))
		offv := bef.Normalize().MulScalar(nozzle)
		newPth = append(newPth, Line{p1.Shift(offv), p2.Shift(offv)})
	}

	closedPath := pth.IsClosed()
	if closedPath {
		newPth = append(newPth, newPth[0]) // add first line to the end
	}
	for i := 0; i < len(newPth)-1; i++ {
		inters := newPth[i].IntersectLine2(newPth[i+1])
		if inters != nil {
			newPth[i].P2 = *inters
			newPth[i+1].P1 = *inters
		} else {
			debug.AddLine(newPth[i], debug.GreenColor)
			debug.AddLine(Line{newPth[i+1].P1, newPth[i+1].P2.Shift(AxisX.MulScalar(15))}, debug.RedColor)
			newPth[i+1].P1 = newPth[i].P2
		}
	}
	if closedPath {
		newPth[0] = newPth[len(newPth)-1]
		newPth = newPth[:len(newPth)-1] //remove added line
	}

	res := Path{}
	res.Points = append(res.Points, newPth[0].P1)
	for i := 0; i < len(newPth); i++ {
		res.Points = append(res.Points, newPth[i].P2)
	}

	return res
}
