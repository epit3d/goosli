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

func CommandsWithTemplates(cmds []gcode.Command, settings Settings) bytes.Buffer{
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}
