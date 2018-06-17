package debug

import (
	"bytes"
	. "github.com/l1va/goosli/primitives"
)

var cfg = Config()

func PointsToDebugFile(ps []Point, filename string) {
	if cfg.Debug {
		var b bytes.Buffer
		for i := 0; i < len(ps)-1; i++ {
			b.WriteString("line ")
			b.WriteString(pointToString(ps[i]))
			b.WriteString(pointToString(ps[i+1]) + "\n")
		}

		ToFile(b, cfg.DebugPath+filename)
	}
}

func pointToString(a Point) string {
	return StrF(a.X) + " " + StrF(a.Y) + " " + StrF(a.Z) + " "
}
