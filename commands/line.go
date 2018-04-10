package commands

import "github.com/l1va/goosli"

type Line struct {
	P1, P2 goosli.Point
}

func(l Line) ToGCode() string{
	return "Line " + l.P1.ToString() +" "+ l.P2.ToString()
}
