package slicers

import (
	"github.com/l1va/goosli/debug"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"log"
)

// SliceByProfile - Slicing on layers by simple algo
func SliceByProfile(mesh *Mesh, settings Settings) gcode.Gcode {
	debug.RecreateFile()
	layers := SliceByVector(mesh, settings.LayerHeight, AxisZ)
	LayersToGcode(layers, "/home/l1va/debug.gcode", settings)

	centers := calculateCenters(layers)
	debug.AddPointsToFile(centers, debug.GreenColor)
	simplified := helpers.SimplifyLine(centers, settings.Epsilon)
	debug.AddPointsToFile(simplified, debug.GreenColor)

	up := mesh
	var down *Mesh
	var gcd gcode.Gcode

	for i := 1; i < len(simplified); i++ {
		v := simplified[i-1].VectorTo(simplified[i])
		if i < len(simplified)-1 {
			var err error
			up, down, err = helpers.CutMesh(up, Plane{simplified[i], AxisZ})
			if err != nil {
				log.Fatal("failed to cut mesh by plane: ", err)
			}
		} else {
			down = up
		}
		angleZ := v.ProjectOnPlane(PlaneXY).Angle(AxisX) + 90
		angleX := v.Angle(AxisZ)

		println("angles: ", angleX, " ", angleZ, "")
		down = down.Rotate(RotationAroundZ(angleZ), OriginPoint)
		down = down.Rotate(RotationAroundX(angleX), OriginPoint)
		gcd.Add(gcode.RotateXZ{angleX, angleZ})

		layers := SliceByVector(down, settings.LayerHeight, AxisZ)
		gcd.Add(gcode.LayersMoving{layers, gcd.LayersCount, settings.PlaneCenterZ})
	}
	return gcd
}
func calculateCenters(layers []Layer) []Point {
	var centers []Point
	for _, layer := range layers {
		centers = append(centers, calculateCenter(layer))
	}
	return centers
}

func calculateCenter(layer Layer) Point {
	x, y, z, count := 0.0, 0.0, 0.0, 0
	for _, path := range layer.Paths {
		crd := FindCentroid(path)
		x += crd.X
		y += crd.Y
		z += crd.Z
		count += 1
	}
	if count > 0 {
		countF := float64(count)
		return Point{x / countF, y / countF, z / countF}
	}
	return OriginPoint
}
