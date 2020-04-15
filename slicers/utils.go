package slicers

import (
	"bytes"
	"github.com/l1va/goosli/gcode"
	. "github.com/l1va/goosli/primitives"
)

func LayersToGcode(layers []Layer, filename string, settings Settings) {
	var b bytes.Buffer
	cmd := gcode.LayersMoving{layers, 0, settings.GetExtrusionParams()}
	cmd.ToGCode(&b)
	ToFile(b, filename)
}
