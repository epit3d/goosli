package commands

import "github.com/l1va/goosli"

type Line struct {
	V1,V2 goosli.Vector
}

func(l Line) ToGCode() string{
	return "TODO"
}
