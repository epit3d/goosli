package slicers

import (
	"bytes"
	"github.com/l1va/goosli/gcode"
	. "github.com/l1va/goosli/primitives"
)

func LayersToGcode(layers []Layer, filename string, settings Settings) {
	var b bytes.Buffer
	gcd := gcode.NewGcode(*settings.GcodeSettings)
	gcd.AddLayers(layers)
	gcd.ToOutput(&b)
	ToFile(b, filename)
}
