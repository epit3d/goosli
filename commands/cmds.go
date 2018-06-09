package commands

import (
	"bytes"
	"math"
	"strconv"
	"github.com/l1va/goosli"
)

type Command interface {
	ToGCode(b *bytes.Buffer)
}

type RotateXZ struct {
	AngleX int
	AngleZ int
}

func (r RotateXZ) ToGCode(b *bytes.Buffer) {
	b.WriteString("G62 " + strconv.Itoa(r.AngleX) + " " + strconv.Itoa(r.AngleZ) + "\n")
}

type LayersMoving struct {
	Layers []goosli.Layer
	Index  int
}

func (lm LayersMoving) ToGCode(b *bytes.Buffer) {
	for i := 0; i < len(lm.Layers); i++ {
		b.WriteString(";Layer" + strconv.Itoa(i+lm.Index) + "\n")
		layerToGCode(lm.Layers[i], b)
	}
}

func layerToGCode(l goosli.Layer, b *bytes.Buffer) {
	eOff := 0.0 //TODO: fix extruder value
	for _, p := range l.Paths {
		b.WriteString("G0 " + pointToString(p.Lines[0].P1) + "\n")
		for _, line := range p.Lines {
			eDist := math.Sqrt(math.Pow(line.P2.X-line.P1.X, 2) + math.Pow(line.P2.Y-line.P1.Y, 2) + math.Pow(line.P2.Z-line.P1.Z, 2))
			eOff += eDist
			b.WriteString("G1 " + pointToString(line.P2) + " E" + goosli.StrF(eOff) + "\n")
		}
	}
}

func pointToString(a goosli.Point) string {
	return "X" + goosli.StrF(a.X) + " Y" + goosli.StrF(a.Y) + " Z" + goosli.StrF(a.Z)
}
