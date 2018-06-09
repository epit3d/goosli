package slicers

import (
	"github.com/l1va/goosli"
	"bytes"
	"github.com/l1va/goosli/commands"
)

func LayersToGcode(layers []goosli.Layer, filename string) {
	var b bytes.Buffer
	cmd := commands.LayersMoving{layers, 0}
	cmd.ToGCode(&b)
	goosli.ToFile(b, filename)
}

func cmdsToBuffer(cmds []commands.Command, b *bytes.Buffer) {
	for _,cmd:= range(cmds){
		cmd.ToGCode(b)
	}
}