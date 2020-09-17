package vip

import (
	"github.com/l1va/goosli/debug"
	. "github.com/l1va/goosli/primitives"
)

func MakeOffset(pth Path, nozzle float64, norm Vector) *Path {
	if len(pth.Points) < 2 {
		return nil
	}
	norm = norm.Reverse()
	//norm := pth.Lines[0].ToVector().Cross(pth.Lines[len(pth.Lines)/2].ToVector()) //TODO: choose normal or have in path

	newPth := []Line{}
	copyNewPth := []Line{}
	for i := 1; i < len(pth.Points); i++ {
		p1 := pth.Points[i-1]
		p2 := pth.Points[i]
		bef := norm.Cross(ToVector(p1, p2))
		offv := bef.Normalize().MulScalar(nozzle)
		newPth = append(newPth, Line{p1.Shift(offv), p2.Shift(offv)})
		copyNewPth = append(copyNewPth, newPth[i-1])
	}

	closedPath := pth.IsClosed() //TODO: refactor everything here
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
	var remove []bool
	var afterRemove []Line
	for i := 0; i < len(newPth); i++ {
		if newPth[i].ToVector().CodirectedWith(copyNewPth[i].ToVector()) {
			remove = append(remove, false)
			afterRemove = append(afterRemove, newPth[i])
		} else {
			remove = append(remove, true) //remove if become reversed
		}
	}

	num := len(afterRemove) - 1
	if closedPath {
		num = len(afterRemove)
	}
	for i := 0; i < num; i++ {
		a := afterRemove[i]
		b := afterRemove[(i+1)%len(afterRemove)]
		if !AlmostZero(a.P2.DistanceTo(b.P1)) {
			inters := a.IntersectLine2(b)
			if inters != nil {
				a.P2 = *inters
				b.P1 = *inters
			} else {
				debug.AddLine(a, debug.BlueColor)
				debug.AddLine(Line{b.P1, b.P2.Shift(AxisX.MulScalar(15))}, debug.BlackColor)
				b.P1 = a.P2
			}
		}
	}

	if len(afterRemove) < 2 {
		return nil
	}

	res := Path{}
	res.Points = append(res.Points, afterRemove[0].P1)
	for i := 0; i < len(afterRemove); i++ {
		res.Points = append(res.Points, afterRemove[i].P2)
	}
	res = res.MinimizeLines()
	return &res
}
