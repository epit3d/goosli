package slicers

import (
	"github.com/l1va/goosli"
	"github.com/l1va/goosli/commands"
)

func Slice3DOF(mesh goosli.Mesh) []commands.Command {

	res := make([]commands.Command, 10)
	for i := 0; i < 10; i++ {
		res[i] = commands.Line{}
	}
	return res
}
