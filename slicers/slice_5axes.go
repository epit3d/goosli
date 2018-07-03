package slicers

import (
	"bytes"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/gcode"
	"log"
	"math"
	"github.com/l1va/goosli/debug"
	"fmt"
)

var
(
	angle_dist  = 10
	angle_full  = 360
	angle_count = angle_full / angle_dist
	angleX = 60.0

)

// SliceByProfile - Slicing on layers by simple algo
func Slice5Axes(mesh *Mesh, settings Settings) bytes.Buffer {
	debug.RecreateFile()
	layersCount := 0
	var cmds []gcode.Command
	simpMesh, err := helpers.SimplifyMesh(mesh, 500)
	if err != nil {
		log.Fatal("simplification during slicing failed:", err)
	}

	Z := AxisZ.Normalize()

	minz, maxz := simpMesh.MinMaxZ(Z)
	n := int(math.Ceil((maxz - minz) / settings.LayerHeight))
	sh := Z.MulScalar(maxz - minz).MulScalar(1.0 / float64(n))
	//curP := Z.MulScalar(minz).ToPoint().Shift(sh.MulScalar(0.5))

	i := 0
	layers := SliceByVector(simpMesh, settings.LayerHeight, AxisZ) //TODO: make slicing by one layer?
	//LayersToGcode(layers, "/home/l1va/debug.gcode")
	toAdd := SliceByVector(mesh, settings.LayerHeight, AxisZ)
	rotated := false
	angleZ := 0.0
	down := mesh
	for i < len(layers) {
		if i == 0 {
			i++
			continue
		}

		newPlane := calculateFails(layers[i], layers[i-1], sh, settings.LayerHeight)

		if newPlane == nil {
			i++
			if i == len(layers) {
				cmds = append(cmds, gcode.LayersMoving{toAdd, layersCount})
				layersCount += len(toAdd)
			}
		} else {
			cmds = append(cmds, gcode.LayersMoving{toAdd[:i], layersCount})
			layersCount += i

			//if rotated {
			//	break
			//}
			//rotated = true
			debug.AddLayer(layers[i-1])
			debug.AddLayer(layers[i])
			curP := calculateCenter(layers[i])
			//fmt.Printf("\n %v \n", layers[i])
			fmt.Printf("\n %v \n", i)

			//SaveSTL("/home/l1va/simp_mesh_before.stl", simpMesh)
			//debug.TriangleToDebugFile(curP, curP.Shift(AxisX.MulScalar(10)), curP.Shift(AxisY.MulScalar(15)), "debug_simplified.txt")
			simpMesh, _, err = helpers.CutMesh(simpMesh, Plane{curP, AxisZ})
			if err != nil {

				log.Fatal("failed to cut simpMesh by plane in 5a slicing: ", err)
			}
			mesh, _, err = helpers.CutMesh(mesh, Plane{curP, AxisZ})
			if err != nil {
				log.Fatal("failed to cut mesh by plane in 5a slicing: ", err)
			}
			debug.AddTriangleByPoints(newPlane.P, newPlane.P.Shift(AxisX.ProjectOnPlane(*newPlane).MulScalar(10)),
				newPlane.P.Shift(AxisY.ProjectOnPlane(*newPlane).MulScalar(15)))

			//newPlane.N.Rotate(RotationAroundZ(90))

			//SaveSTL("/home/l1va/simp_mesh_before_cut.stl", simpMesh)

			simpMesh, _, err = helpers.CutMesh(simpMesh, *newPlane)
			if err != nil {
				break //TODO: remvoe me
				log.Fatal("failed to cut simpMesh by newPlane in 5a slicing: ", err)
			}
			if i > 10 {
				debug.AddTriangleByPoints(newPlane.P, newPlane.P.Shift(AxisX.ProjectOnPlane(*newPlane).MulScalar(10)),
					newPlane.P.Shift(AxisY.ProjectOnPlane(*newPlane).MulScalar(15)))

				//newPlane.N.Rotate(RotationAroundZ(90))

				//SaveSTL("/home/l1va/simp_mesh_before_cut.stl", simpMesh)

			}
			mesh, down, err = helpers.CutMesh(mesh, *newPlane)
			if err != nil {
				log.Fatal("failed to cut mesh by newPlane in 5a slicing: ", err)
			}
			toAdd = SliceByVector(down, settings.LayerHeight, AxisZ)
			cmds = append(cmds, gcode.LayersMoving{toAdd, layersCount})
			layersCount += len(toAdd)

			if rotated {
				simpMesh = simpMesh.Rotate(RotationAroundX(-angleX), OriginPoint)
				simpMesh = simpMesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)
				mesh = mesh.Rotate(RotationAroundX(-angleX), OriginPoint)
				mesh = mesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)
				cmds = append(cmds, gcode.RotateXZ{0, 0})
				rotated = false
			} else {
				angleZ = newPlane.N.ProjectOnPlane(PlaneXY).Angle(AxisX) + 90
				print(angleZ)
				simpMesh = simpMesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
				simpMesh = simpMesh.Rotate(RotationAroundX(angleX), OriginPoint)
				mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
				mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)
				cmds = append(cmds, gcode.RotateXZ{angleX, angleZ})
				rotated = true
			}
			layers = SliceByVector(simpMesh, settings.LayerHeight, AxisZ)
			toAdd = SliceByVector(mesh, settings.LayerHeight, AxisZ)
			//cmds = append(cmds, gcode.LayersMoving{layers, layersCount})
			//break
			i = 0
			minz, maxz := simpMesh.MinMaxZ(Z)
			n := int(math.Ceil((maxz - minz) / settings.LayerHeight))
			sh = Z.MulScalar(maxz - minz).MulScalar(1.0 / float64(n))
		}
	}

	settings.LayerCount = layersCount
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}

func calculateFails(layer, prev_layer Layer, sh Vector, layerHeight float64) *Plane {
	half_height := layerHeight + 0.1 //TODO: remove hardcode
	rm := RotationAroundZ(float64(angle_dist))

	for _, path := range (layer.Paths) {
		curP := calculateCenterForPath(path)
		prevP := curP.Shift(sh.Reverse())
		prev_path := findPrevPath(prevP, prev_layer)
		failsCur := []float64{}
		if prev_path == nil {
			prevPath := prev_layer.Paths[0] //checking only first path
			pathP := calculateCenterForPath(prevPath)
			for _, line := range (prevPath.Lines) { //TODO: optimize me
				pl := Plane{pathP, pathP.VectorTo(prevP).Cross(AxisZ)}
				pi := pl.IntersectSegment(line.P1, line.P2)
				if pi != nil && pathP.VectorTo(*pi).CodirectedWith(pathP.VectorTo(prevP)) {
					angle := pathP.VectorTo(*pi).Angle(AxisX)+90
					norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angle))
					return &Plane{*pi, norm}
				}
			}
		} else {
			curX := AxisX
			for j := 0; j < angle_count; j++ {
				cur_dist := 0.0
				for _, line := range (path.Lines) { //TODO: optimize me
					pl := Plane{curP, curX.Cross(AxisZ)}
					pi := pl.IntersectSegment(line.P1, line.P2)
					if pi != nil && curP.VectorTo(*pi).CodirectedWith(curX) {
						cur_dist = curP.DistanceTo(*pi)
						break
					}
				}
				prev_dist := 0.0
				for _, line := range (prev_path.Lines) { //TODO: optimize me
					pl := Plane{prevP, curX.Cross(AxisZ)}
					pi := pl.IntersectSegment(line.P1, line.P2)
					if pi != nil && prevP.VectorTo(*pi).CodirectedWith(curX) {
						prev_dist = prevP.DistanceTo(*pi)
						break
					}
				}
				curX = curX.Rotate(rm)
				if prev_dist == 0 || cur_dist == 0 {
					continue
				}
				if cur_dist > prev_dist+half_height {
					failsCur = append(failsCur, float64(j*angle_dist))
				}
			}
			if len(failsCur) > 0 {
				prevPath := *prev_path
				pathP := calculateCenterForPath(prevPath)
				angle := float64(MiddleOnTheRing(failsCur, angle_full))
				v := AxisX.Rotate(RotationAroundZ(angle))
				pl := Plane{pathP, v.Cross(AxisZ)}

				for _, line := range prevPath.Lines { //TODO: optimize me
					pi := pl.IntersectSegment(line.P1, line.P2)
					if pi != nil && pathP.VectorTo(*pi).CodirectedWith(v) {
						norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angle + 90))
						return &Plane{*pi, norm}
					}
				}
			}
		}

	}

	return nil

}

func findPrevPath(prevP Point, prev_layer Layer) *Path {
	for _, path := range (prev_layer.Paths) {
		if prevP.Inside(path) {
			return &path
		}
	}
	return nil
}

func MiddleOnTheRing(arr []float64, biggest int) int {
	// assume that arr is sorted
	i := 0
	for arr[len(arr)-1]-arr[i] > float64(biggest)/2 {
		arr[i] += float64(biggest)
		i++
	}
	sum := 0.0
	for _, num := range arr {
		sum += num
	}
	return int(sum/float64(len(arr))) % biggest
}
