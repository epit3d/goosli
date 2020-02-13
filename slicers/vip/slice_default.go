package vip

import (
	"github.com/l1va/goosli/debug"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	. "github.com/l1va/goosli/slicers"
	"log"
)

func SliceDefault2(mesh *Mesh, settings Settings, layers []Layer) gcode.Gcode {

	print("DEFAULT slicing:\n")
	debug.RecreateFile()

	var gcd gcode.Gcode
	fillPlanes := CalcFillPlanes(mesh, settings)

	i := 1
	rotated := false
	angleZ := 0.0
	var down *Mesh
	var err error

	//layes := []Layer{}
	//for _, t := range mesh.Triangles {
	//	lay := Layer{Paths: []Path{Path{Lines: []Line{Line{t.P1, t.P1.Shift(t.N .Normalize())}}}}}
	//	layes = append(layes, lay)
	//}
	//
	//cmds = append(cmds, gcode.LayersMoving{Layers: layes, Index: layersCount})
	//settings.LayerCount = 5
	//return CommandsWithTemplates(cmds, settings)

	for i < len(layers) {

		newPlane := CalculateFails(layers[i-1], layers[i])

		if newPlane == nil {
			i++
			if i == len(layers) {
				gcd.Add(gcode.LayersMoving{Layers: PrepareLayers(layers, settings, fillPlanes), Index: gcd.LayersCount})
			}
		} else {
			mesh, down, err = helpers.CutMesh(mesh, *newPlane)
			if err != nil {
				log.Fatal("failed to cut mesh by newPlane in default slicing: ", err)
			}
			add := SliceByVector(down, settings.LayerHeight, AxisZ)

			gcd.Add(gcode.LayersMoving{Layers: PrepareLayers(add, settings, fillPlanes), Index: gcd.LayersCount})
			println("added: ", len(add), i)
			//debug.AddLayer(layers[i])
			if rotated {
				mesh = mesh.Rotate(RotationAroundX(-angleX), OriginPoint)
				mesh = mesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)

				gcd.Add(gcode.InclineXBack{})
				gcd.Add(gcode.RotateZ{Angle: 0})
				rotated = false
				println("triangles: ", len(mesh.Triangles))
				println("rotate back")
				for j, plane := range fillPlanes {
					fillPlanes[j] = plane.Rotate(RotationAroundX(-angleX)).Rotate(RotationAroundZ(-angleZ))
				}
			} else {
				angleZ = newPlane.N.ProjectOnPlane(PlaneXY).Angle(AxisX) + 90
				println(angleZ)
				mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
				mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)

				gcd.Add(gcode.RotateZ{Angle: angleZ})
				gcd.Add(gcode.InclineX{})
				rotated = true
				for j, plane := range fillPlanes {
					fillPlanes[j] = plane.Rotate(RotationAroundZ(angleZ)).Rotate(RotationAroundX(angleX))
				}
			}
			layers = SliceByVector(mesh, settings.LayerHeight, AxisZ)
			println("len layers: ", len(layers))

			i = 1
		}
	}
	return gcd

}

/*
func SliceDefault(mesh *Mesh, settings Settings, layers []Layer) bytes.Buffer {

	print("DEFAULT slicing:\n")
	debug.RecreateFile()

	var cmds []gcode.Command
	layersCount := 0
	rotated := false
	for {
		centers := calculateCenters(layers)
		simplified := helpers.SimplifyLine(centers, settings.Epsilon)

		debug.AddPointsToFile(simplified)
		j := 1
		for j < len(simplified) {
			println(AxisZ.Angle(simplified[j-1].VectorTo(simplified[j])))
			if AxisZ.Angle(simplified[j-1].VectorTo(simplified[j])) > failAngle {
				break
			}
			j++
		}
		if j == len(simplified) {
			cmds = append(cmds, gcode.LayersMoving{layers, layersCount})
			layersCount += len(layers)
			break
		}
		projectedV := simplified[j-1].VectorTo(simplified[j]).ProjectOnPlane(PlaneXY)
		angleZ := projectedV.Angle(AxisX) + 90
		norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angleZ))
		println("Z: ", angleZ)
		debug.AddLine(Line{simplified[j-1], simplified[j-1].Shift(norm.MulScalar(10))})
		// trying to find point for newPlane
		planeHor := Plane{simplified[j-1], AxisZ}
		layer := planeHor.IntersectMesh(mesh)
		debug.AddLayer(layer)
		if len(layer.Paths) > 1 {
			println("len path is not 1: ", len(layer.Paths))
		}
		planeVert := Plane{simplified[j-1], projectedV.Cross(AxisZ)}
		p := planeVert.IntersectPathCodirectedWith(layer.Paths[0], projectedV)
		if p != nil {
		} else {
			println("point was not found, using center: ")
			p = &simplified[j-1]
		}
		// cutting and slicing and rotating
		newPlane := Plane{*p, norm}
		var err error
		var down *Mesh
		mesh, down, err = helpers.CutMesh(mesh, newPlane)
		if err != nil {
			log.Fatal("failed to cut mesh by plane: ", err)
		}
		add := SliceByVector(down, settings.LayerHeight, AxisZ)
		cmds = append(cmds, gcode.LayersMoving{add, layersCount})
		layersCount += len(add)

		if rotated {
			mesh = mesh.Rotate(RotationAroundX(-angleX), OriginPoint)
			mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)

			cmds = append(cmds, gcode.RotateXZ{0, angleZ})
			rotated = false
		} else {
			mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
			mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)

			cmds = append(cmds, gcode.RotateXZ{angleX, angleZ})
			rotated = true
		}

		layers = SliceByVector(mesh, settings.LayerHeight, AxisZ)
	}

	settings.LayerCount = layersCount
	return CommandsWithTemplates(cmds, settings)

}*/

func calculateCenters(layers []Layer) []Point { //TODO:
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
