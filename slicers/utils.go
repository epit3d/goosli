package slicers

import (
	. "github.com/l1va/goosli/primitives"
	"bytes"
	"github.com/l1va/goosli/gcode"
)

func LayersToGcode(layers []Layer, filename string) {
	var b bytes.Buffer
	cmd := gcode.LayersMoving{layers, 0}
	cmd.ToGCode(&b)
	ToFile(b, filename)
}

func cmdsToBuffer(cmds []gcode.Command, b *bytes.Buffer) {
	for _, cmd := range (cmds) {
		cmd.ToGCode(b)
	}
}
