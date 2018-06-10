package debug

import (
	"bytes"
	"github.com/l1va/goosli"
)

var cfg = Config()

func PointsToDebugFile(ps []goosli.Point, filename string) {
	if cfg.Debug {
		var b bytes.Buffer
		for i := 0; i < len(ps)-1; i++ {
			b.WriteString("line ")
			b.WriteString(pointToString(ps[i]))
			b.WriteString(pointToString(ps[i+1]) + "\n")
		}

		goosli.ToFile(b, cfg.DebugPath+filename)
	}
}

func pointToString(a goosli.Point) string {
	return goosli.StrF(a.X) + " " + goosli.StrF(a.Y) + " " + goosli.StrF(a.Z) + " "
}
