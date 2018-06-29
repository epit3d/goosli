package slicers

import (
	"bytes"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/gcode"
	"log"
	"math"
	"fmt"
	"github.com/l1va/goosli/debug"
)

var
(
	angle_dist  = 10
	angle_full  = 360
	angle_count = angle_full / angle_dist
)

// SliceByProfile - Slicing on layers by simple algo
func Slice5Axes(mesh *Mesh, settings Settings) bytes.Buffer {

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
	LayersToGcode(layers, "/home/l1va/debug.gcode")
	toAdd := SliceByVector(mesh, settings.LayerHeight, AxisZ)
	//rotated := false
	angleX := 0.0
	angleZ := 0.0
	for i < len(layers) {
		if i == 0 {
			i++
			continue
		}

		fails := calculateFails(layers[i], layers[i-1], sh, settings.LayerHeight)

		if len(fails) == 0 {
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
			SaveSTL("/home/l1va/simp_mesh_before_cut.stl", simpMesh)
			curP := calculateCenter(layers[i])
			if i != len(layers)-1 {
				fmt.Printf("\n %v \n", layers[i])
				fmt.Printf("\n %v \n", layers[i-1])
				debug.TriangleToDebugFile(curP, curP.Shift(AxisX.MulScalar(10)), curP.Shift(AxisY.MulScalar(15)), "debug_simplified.txt")
				simpMesh, _, err = helpers.CutMesh(simpMesh, Plane{curP, AxisZ})
				if err != nil {
					log.Fatal("failed to cut simpMesh by plane in 5a slicing: ", err)
				}
				mesh, _, err = helpers.CutMesh(mesh, Plane{curP, AxisZ})
				if err != nil {
					log.Fatal("failed to cut mesh by plane in 5a slicing: ", err)
				}
			}
			SaveSTL("/home/l1va/simp_mesh_before.stl", simpMesh)
			//revert previous to zero
			simpMesh = simpMesh.Rotate(RotationAroundX(-angleX), OriginPoint)
			simpMesh = simpMesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)

			layers = SliceByVector(simpMesh, settings.LayerHeight, AxisZ)
			LayersToGcode(layers, "/home/l1va/debug2.gcode")
			SaveSTL("/home/l1va/simp_mesh.stl", simpMesh)
			curX := 0.0
			curZ := 0.0
			if len(layers) > 1 {
				fails = calculateFails(layers[1], layers[0], sh, settings.LayerHeight)
				fmt.Printf("\n%v first compare: %v %v", fails, len(layers[0].Paths), len(layers[1].Paths))
				fmt.Printf("%v", layers[0])
				for _, fail := range fails {
					curX = 60
					curZ = float64(fail) + 90
					simpMesh = simpMesh.Rotate(RotationAroundZ(curZ), OriginPoint)
					simpMesh = simpMesh.Rotate(RotationAroundX(curX), OriginPoint)
					layers = SliceByVector(simpMesh, settings.LayerHeight, AxisZ)
					if len(layers) > 1 {
						fails = calculateFails(layers[1], layers[0], sh, settings.LayerHeight)
						if len(fails) == 0 {
							break
						}
						simpMesh = simpMesh.Rotate(RotationAroundX(-curX), OriginPoint)
						simpMesh = simpMesh.Rotate(RotationAroundZ(-curZ), OriginPoint)
					}

				}
			}
			mesh = mesh.Rotate(RotationAroundX(-angleX), OriginPoint)
			mesh = mesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)

			angleX = curX
			angleZ = curZ
			print("\n !!!!!!!!!! MAXZ: ", angleX, angleZ)

			cmds = append(cmds, gcode.RotateXZ{angleX, angleZ})

			mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
			mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)
			toAdd = SliceByVector(mesh, settings.LayerHeight, AxisZ)
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

func calculateFails(layer, prev_layer Layer, sh Vector, layerHeight float64) []int {
	half_height := layerHeight + 0.1 //TODO: remove hardcode
	rm := RotationAroundZ(float64(angle_dist))
	fails := []int{}

	for _, path := range (layer.Paths) {
		curP := calculateCenterForPath(path)
		prevP := curP.Shift(sh.Reverse())
		prev_path := findPrevPath(prevP, prev_layer)
		failsCur := []float64{}
		if prev_path == nil {
			for _, prev_path := range (prev_layer.Paths) {
				pathP := calculateCenterForPath(prev_path)
				angle := pathP.VectorTo(prevP).Angle(AxisX)
				failsCur = append(failsCur, angle)
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
		}
		if len(failsCur) > 0 {
			fails = append(fails, MiddleOnTheRing(failsCur, angle_full))
		}
	}

	if len(fails) > 0 {
		fmt.Printf("\n %v\n", fails)
		/*		print("\n======= layer ")
				for j := 0; j < len(fails); j++ {
					print("\n", fails[j], " ", cur_dist[fail[j]], " ", prev_dist[fail[j]])
					if !rotated {
						points = append(points, prevP)
						points = append(points, prev_p[fail[j]])
						points1 = append(points1, curP)
						points1 = append(points1, cur_p[fail[j]])
					}
				}
				debug.PointsToDebugFile(points, "debug.txt")
				debug.PointsToDebugFile(points1, "debug_simplified.txt")*/
		print("\n---------")
	}
	return fails
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
