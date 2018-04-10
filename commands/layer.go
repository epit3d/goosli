package commands

import (
	"bytes"
	"math"
	"github.com/l1va/goosli"
)

type Layer struct {
	Order float64
	Paths []Path
}

func (l Layer) ToGCode() string {
	var buffer bytes.Buffer
	buffer.WriteString(";Layer\n")
	eOff := 0.0
	for _, p := range l.Paths {
		buffer.WriteString("G0 " + p.Lines[0].P1.ToString() + "\n")
		for _, line := range p.Lines {
			eDist := math.Sqrt(math.Pow(line.P2.X-line.P1.X, 2) + math.Pow(line.P2.Y-line.P1.Y, 2) + math.Pow(line.P2.Z-line.P1.Z, 2))
			buffer.WriteString("G1 " + line.P2.ToString() + " E" + goosli.StrF(eOff+eDist) + "\n")
			eOff += eDist
		}
	}
	return buffer.String()
}
