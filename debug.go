package goosli

import "bytes"

func PointsToDebugFile(ps []Point, filename string) {
	var b bytes.Buffer
	for i := 0; i < len(ps)-1; i++ {
		b.WriteString("line ")
		b.WriteString(pointToString(ps[i]))
		b.WriteString(pointToString(ps[i+1]) + "\n")
	}

	ToFile(b, filename)
}

func pointToString(a Point) string {
	return StrF(a.X) + " " + StrF(a.Y) + " " + StrF(a.Z) + " "
}
