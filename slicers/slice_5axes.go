package slicers

import (
	"bytes"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/gcode"
	"log"
	"math"
	"github.com/l1va/goosli/debug"
)

// SliceByProfile - Slicing on layers by simple algo
func Slice5Axes(mesh *Mesh, settings Settings) bytes.Buffer {

	layersCount := 0
	var cmds []gcode.Command
	mesh, err := helpers.SimplifyMesh(mesh, 500)
	if err != nil {
		log.Fatal("simplification during slicing failed:", err)
	}

	angle_dist := 10
	angle_count := 360 / angle_dist

	Z := AxisZ.Normalize()

	minz, maxz := mesh.MinMaxZ(Z)
	n := int(math.Ceil((maxz - minz) / settings.LayerHeight))
	sh := Z.MulScalar(maxz - minz).MulScalar(1.0 / float64(n))
	//curP := Z.MulScalar(minz).ToPoint().Shift(sh.MulScalar(0.5))

	prev_dist := [] float64{}
	cur_dist := [] float64{}
	prev_p := []Point{}
	cur_p := []Point{}
	X := AxisX.Normalize()

	for i := 0; i < angle_count; i++ {
		prev_dist = append(prev_dist, 0)
		cur_dist = append(cur_dist, 0)
		prev_p = append(prev_p, OriginPoint)
		cur_p = append(cur_p, OriginPoint)
	}

	half_height := settings.LayerHeight + 0.1
	rm := RotationAroundZ(float64(angle_dist))
	X = AxisX
	i := 0
	layers := SliceByVector(mesh, settings.LayerHeight, AxisZ) //TODO: make slicing by one layer
	LayersToGcode(layers, "/home/l1va/debug.gcode")
	points := [] Point{}
	points1 := [] Point{}
	rotated := false
	for i < len(layers) {
		if i == 0 {
			i++
			continue
		}
		layer := layers[i]
		prev_layer := layers[i-1]
		curP := calculateCenter(layer)
		prevP := curP.Shift(sh.Reverse())
		curX := X
		fail := []int{}

		for j := 0; j < angle_count; j++ {
			path := layer.Paths[0] //TODO: refactor if not one path
			if len(layer.Paths) > 1 {
				//print("\nWarning!!! len paths is ", len(layer.Paths))
			}
			cur_dist[j] = 0
			for _, line := range (path.Lines) { //TODO: optimize me
				pl := Plane{curP, curX.Cross(Z)}
				pi := pl.IntersectSegment(line.P1, line.P2)
				if pi != nil && curP.VectorTo(*pi).CodirectedWith(curX) {
					cur_dist[j] = curP.DistanceTo(*pi)
					cur_p[j] = *pi
					break
				}
			}
			prevpath := prev_layer.Paths[0] //TODO: refactor if not one path
			if len(prev_layer.Paths) > 1 {
				//print("\nWarning!!! len paths is ", len(layer.Paths))
			}
			prev_dist[j] = 0
			for _, line := range (prevpath.Lines) { //TODO: optimize me
				pl := Plane{prevP, curX.Cross(Z)}
				pi := pl.IntersectSegment(line.P1, line.P2)
				if pi != nil && prevP.VectorTo(*pi).CodirectedWith(curX) {
					prev_dist[j] = prevP.DistanceTo(*pi)
					prev_p[j] = *pi
					break
				}
			}
			curX = curX.Rotate(rm)
			if prev_dist[j] == 0 || cur_dist[j] == 0 {
				continue
			}
			if cur_dist[j] > prev_dist[j]+half_height {
				fail = append(fail, j)
			}
		}

		if len(fail) == 0 {
			i++
			if i == len(layers) {
				cmds = append(cmds, gcode.LayersMoving{layers, layersCount})
			}
		} else {
			print("\n", len(fail))
			print("\n======= layer ", i+layersCount)
				for j := 0; j < len(fail); j++ {
					print("\n", fail[j]*angle_dist, " ", cur_dist[fail[j]], " ", prev_dist[fail[j]])
					if !rotated {
						points = append(points, prevP)
						points = append(points, prev_p[fail[j]])
						points1 = append(points1, curP)
						points1 = append(points1, cur_p[fail[j]])
					}
				}
				print("---------")

			angleX := 60.0
			angleZ := float64(middleOnTheRing(fail, angle_count) * angle_dist ) + 90
			print("\n !!!!!!!!!! MAXZ: ", angleZ)

			cmds = append(cmds, gcode.LayersMoving{layers[:i], layersCount})
			if rotated {
				break
			}
			rotated = true
			cmds = append(cmds, gcode.RotateXZ{angleX, angleZ})
			layersCount += i

			if i != len(layers)-1 {
				mesh, _, err = helpers.CutMesh(mesh, Plane{curP, AxisZ})
				if err != nil {
					log.Fatal("failed to cut mesh by plane in 5a slicing: ", err)
				}
			}
			mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
			mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)
			layers = SliceByVector(mesh, settings.LayerHeight, AxisZ)
			i = 0
			minz, maxz := mesh.MinMaxZ(Z)
			n := int(math.Ceil((maxz - minz) / settings.LayerHeight))
			sh = Z.MulScalar(maxz - minz).MulScalar(1.0 / float64(n))
		}
	}
	debug.PointsToDebugFile(points, "debug.txt")
	debug.PointsToDebugFile(points1, "debug_simplified.txt")
	settings.LayerCount = layersCount
	smap := settings.ToMap()

	var buffer bytes.Buffer
	buffer.WriteString(PrepareDataFile("data/header_template.txt", smap))
	cmdsToBuffer(cmds, &buffer)
	buffer.WriteString(PrepareDataFile("data/footer_template.txt", smap))
	return buffer
}
func middleOnTheRing(arr []int, biggest int) int {
	// assume that arr is sorted
	i := 0
	for arr[len(arr)-1]-arr[i] > biggest/2 {
		arr[i] += biggest
		i++
	}
	sum := 0
	for _, num := range arr {
		sum += num
	}
	return (sum / len(arr)) % biggest
}
