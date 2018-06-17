package slicers

import (
	. "github.com/l1va/goosli/primitives"
	"bytes"
	"log"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	"github.com/l1va/goosli/debug"
)

// SliceByProfile - Slicing on layers by simple algo
func SliceByProfile(mesh *Mesh, epsilon float64, settings Settings) bytes.Buffer {
	layers := SliceByVector(mesh, settings.LayerHeight, AxisZ)
	//LayersToGcode(layers, "debug.gcode")

	centers := calculateCenters(layers)
	debug.PointsToDebugFile(centers, "debug.txt")

	simplified := helpers.SimplifyLine(centers, epsilon)
	debug.PointsToDebugFile(simplified, "debug_simplified.txt")

	layersCount := 0
	up := mesh
	var down *Mesh
	var cmds []gcode.Command

	for i := 1; i < len(simplified); i++ {
		v := simplified[i-1].VectorTo(simplified[i])
		if i < len(simplified)-1 {
			var err error
			up, down, err = helpers.CutMesh(up, Plane{simplified[i], v})
			if err != nil {
				log.Fatal("failed to cut mesh by plane: ", err)
			}
		} else {
			down = up
		}
		angleZ := v.ProjectOnPlane(PlaneXY).Angle(AxisX)
		angleX := v.ProjectOnPlane(PlaneYZ).Angle(AxisZ)

		down = down.Rotate(RotationAroundX(angleX), OriginPoint)
		down = down.Rotate(RotationAroundZ(angleZ), OriginPoint) // local rotation!!!
		cmds = append(cmds, gcode.RotateXZ{angleX, angleZ})

		layers := SliceByVector(down, settings.LayerHeight, AxisZ)
		cmds = append(cmds, gcode.LayersMoving{layers, layersCount})
		layersCount += len(layers)
	}
	settings.LayerCount = layersCount
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}

func calculateCenters(layers []Layer) []Point {
	var centers []Point
	for _, layer := range layers {
		x, y, z, count := 0.0, 0.0, 0.0, 0
		for _, path := range layer.Paths {
			for _, line := range path.Lines {
				x += line.P1.X + line.P2.X
				y += line.P1.Y + line.P2.Y
				z += line.P1.Z + line.P2.Z
			}
			count += len(path.Lines) * 2
		}
		if count > 0 {
			countF := float64(count)
			centers = append(centers, Point{x / countF, y / countF, z / countF})
		}
	}
	return centers
}