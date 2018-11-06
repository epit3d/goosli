package debug

import (
	"bytes"
	. "github.com/l1va/goosli/primitives"
	"os"
	"fmt"
)

type DebugColor string

var (
	RedColor   = DebugColor("Red")
	BlueColor  = DebugColor("Blue")
	GreenColor = DebugColor("Green")
	BlackColor = DebugColor("Black")
)
var cfg = Config()

func RecreateFile() {
	// delete file
	var _, err = os.Stat(cfg.DebugFile)

	// create file if not exists
	if os.IsExist(err) {
		var err = os.Remove(cfg.DebugFile)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	var file *os.File
	file, err = os.Create(cfg.DebugFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
}

func AddPointsToFile(ps []Point, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		for i := 0; i < len(ps)-1; i++ {
			b.WriteString("line "+ string(color)+" ")
			b.WriteString(pointToString(ps[i]))
			b.WriteString(pointToString(ps[i+1]) + "\n")
		}
		AddToFile(b, cfg.DebugFile)
	}
}

func AddLine(l Line, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		b.WriteString("line "+ string(color)+" ")
		b.WriteString(pointToString(l.P1))
		b.WriteString(pointToString(l.P2) + "\n")

		AddToFile(b, cfg.DebugFile)
	}
}

func AddTriangle(t Triangle, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		b.WriteString("triangle "+ string(color)+" ")
		b.WriteString(pointToString(t.P1))
		b.WriteString(pointToString(t.P2))
		b.WriteString(pointToString(t.P3) + "\n")

		AddToFile(b, cfg.DebugFile)
	}
}
func AddTriangleByPoints(p1, p2, p3 Point, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		b.WriteString("triangle "+ string(color)+" ")
		b.WriteString(pointToString(p1))
		b.WriteString(pointToString(p2))
		b.WriteString(pointToString(p3) + "\n")

		AddToFile(b, cfg.DebugFile)
	}
}

func AddLayer(layer Layer, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		for _, path := range layer.Paths {
			for _, line := range path.Lines {
				b.WriteString("line "+ string(color)+" ")
				b.WriteString(pointToString(line.P1))
				b.WriteString(pointToString(line.P2) + "\n")
			}
		}
		AddToFile(b, cfg.DebugFile)
	}
}

func AddPath(path Path, color DebugColor) {
	if cfg.Debug {
		var b bytes.Buffer
		for _, line := range path.Lines {
			b.WriteString("line "+ string(color)+" ")
			b.WriteString(pointToString(line.P1))
			b.WriteString(pointToString(line.P2) + "\n")
		}
		AddToFile(b, cfg.DebugFile)
	}
}

func pointToString(a Point) string {
	return StrF(a.X) + " " + StrF(a.Y) + " " + StrF(a.Z) + " "
}
