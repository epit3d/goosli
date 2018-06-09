package slicers

import (
	"github.com/l1va/goosli"
	"bytes"
	"strconv"
	"log"
)

// SliceByProfile - Slicing on layers by simple algo
func SliceByProfile(mesh *goosli.Mesh, thickness float64, epsilon float64) bytes.Buffer {

	layers := SliceByVector(mesh, thickness, goosli.AxisZ)
	//goosli.LayersToGcode(layers, "/home/l1va/debug.gcode")

	centers := calculateCenters(layers)
	//goosli.PointsToDebugFile(centers, "/home/l1va/debug.txt")

	simplified := SimplifyLine(centers, epsilon)
	//goosli.PointsToDebugFile(simplified, "/home/l1va/debug_simplified.txt")

	var buffer bytes.Buffer
	layersCount := 0
	up := mesh
	var down *goosli.Mesh

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

		down = rotateXZ(angleX, angleZ, down, &buffer, goosli.OriginPoint)

		layers := SliceByVector(mesh, thickness, v)
		goosli.LayersToBuffer(layers, layersCount, &buffer)
		layersCount += len(layers)
	}

	return buffer
}



func rotateXZ(angleX int, angleZ int, mesh *goosli.Mesh, b *bytes.Buffer, around goosli.Point) *goosli.Mesh {
	b.WriteString("G62 " + strconv.Itoa(angleX) + " " + strconv.Itoa(angleZ) + "\n")
	mesh = mesh.RotateX(angleX, around)
	return mesh.RotateZ(angleZ, around) // local rotation!!!
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
