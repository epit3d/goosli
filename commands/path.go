package commands

import (
	"bytes"
	"github.com/l1va/goosli"
)

type Path struct {
	Lines []Line
}

func (p Path) ToGCode() string {
	var buffer bytes.Buffer
	buffer.WriteString("Path gCode\n")
	for _, l := range p.Lines {
		buffer.WriteString(l.ToGCode() + "\n")
	}
	return buffer.String()
}

func JoinPaths(paths []Path) []Path {
	lookup := make(map[goosli.Point]Path, len(paths))
	for _, path := range paths {
		lookup[path.Lines[0].P1] = path
	}
	var result []Path
	for len(lookup) > 0 {
		var v goosli.Point
		for v = range lookup {
			break
		}
		var path Path
		for {
			if p, ok := lookup[v]; ok {
				path.Lines = append(path.Lines, p.Lines[0])
				delete(lookup, v)
				v = p.Lines[len(p.Lines)-1].P2
			} else {
				break
			}
		}
		result = append(result, path)
	}
	return result
}
