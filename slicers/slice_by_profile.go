package slicers

import (
	"github.com/l1va/goosli"
	"bytes"
	"log"
	"github.com/l1va/goosli/commands"
	"github.com/l1va/goosli/debug"
)

// SliceByProfile - Slicing on layers by simple algo
func SliceByProfile(mesh *goosli.Mesh, epsilon float64, settings Settings) bytes.Buffer {
	layers := SliceByVector(mesh, settings.LayerHeight, goosli.AxisZ)
	//LayersToGcode(layers, "debug.gcode")

	centers := calculateCenters(layers)
	debug.PointsToDebugFile(centers, "debug.txt")

	simplified := SimplifyLine(centers, epsilon)
	debug.PointsToDebugFile(simplified, "debug_simplified.txt")

	layersCount := 0
	up := mesh
	var down *goosli.Mesh
	var cmds []commands.Command

	for i := 1; i < len(simplified); i++ {
		v := simplified[i-1].VectorTo(simplified[i])
		if i < len(simplified)-1 {
			var err error
			up, down, err = BisectMesh(up, goosli.Plane{simplified[i], v})
			if err != nil {
				log.Fatal("failed to cut mesh by plane: ", err)
			}
		} else {
			down = up
		}
		angleZ := int(v.ProjectOnPlane(goosli.PlaneXY).Angle(goosli.AxisX))
		angleX := int(v.ProjectOnPlane(goosli.PlaneYZ).Angle(goosli.AxisZ))

		down = down.RotateX(angleX, goosli.OriginPoint)
		down = down.RotateZ(angleZ, goosli.OriginPoint) // local rotation!!!
		cmds = append(cmds, commands.RotateXZ{angleX, angleZ})

		layers := SliceByVector(down, settings.LayerHeight, v)  //TODO: should be sliced before or v rotated
		cmds = append(cmds, commands.LayersMoving{layers, layersCount})
		layersCount += len(layers)
	}
	settings.LayerCount = layersCount
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(goosli.PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(goosli.PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}

func calculateCenters(layers []goosli.Layer) []goosli.Point {
	var centers []goosli.Point
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
			centers = append(centers, goosli.Point{x / countF, y / countF, z / countF})
		}
	}
	return centers
}
