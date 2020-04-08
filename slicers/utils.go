package slicers

import (
	"bytes"
	"github.com/l1va/goosli/gcode"
	. "github.com/l1va/goosli/primitives"
)

func LayersToGcode(layers []Layer, filename string, settings Settings) {
	var b bytes.Buffer
	extParams := ExtrusionParams{settings.BarDiameter, settings.Flow, settings.LayerHeight, settings.LineWidth}
	cmd := gcode.LayersMoving{layers, 0, extParams}
	cmd.ToGCode(&b)
	ToFile(b, filename)
}
