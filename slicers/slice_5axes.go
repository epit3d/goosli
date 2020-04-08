package slicers

import (
	"fmt"
	"github.com/l1va/goosli/debug"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"log"
)

var (
	angleDist  = 10
	angleFull  = 360
	angleCount = angleFull / angleDist
	angleX     = 60.0
	failAngle  = 40.0
)

// SliceByProfile - Slicing on layers by simple algo
func Slice5Axes(mesh *Mesh, settings Settings) gcode.Gcode {
	debug.RecreateFile()
	var gcd gcode.Gcode

	simpMesh, err := helpers.SimplifyMesh(mesh, 2500) //TODO: hardcoded value
	if err != nil {
		log.Fatal("simplification during slicing failed:", err)
	}

	//t1 := simpMesh.Triangles[199]

	/*	debug.AddLine(Line{t1.P1, t1.P2})
		debug.AddLine(Line{t1.P2, t1.P3})
		debug.AddLine(Line{t1.P3, t1.P1})
		debug.AddLine(Line{t1.P1, t1.P1.Shift(t1.N.MulScalar(10))})*/
	//curP := Z.MulScalar(minz).ToPoint().Shift(sh.MulScalar(0.5))

	i := 0
	layers := SliceByVector(simpMesh, AxisZ, settings) //TODO: make slicing by one layer?
	//LayersToGcode(layers, "/home/l1va/debug.gcode")
	toAdd := SliceByVector(mesh, AxisZ, settings)
	rotated := false
	angleZ := 0.0
	down := mesh
	for i < len(layers) {
		if i == 0 {
			i++
			continue
		}

		newPlane := CalculateFails(layers[i-1], layers[i])
		extParams := ExtrusionParams{settings.BarDiameter, settings.Flow, settings.LayerHeight, settings.LineWidth}
		if newPlane == nil {
			i++
			if i == len(layers) {
				gcd.Add(gcode.LayersMoving{Layers: toAdd, Index: gcd.LayersCount, ExtParams: extParams})
			}
		} else {
			gcd.Add(gcode.LayersMoving{Layers: toAdd[:i], Index: gcd.LayersCount, ExtParams: extParams})

			//if rotated {
			//	break
			//}
			//rotated = true
			debug.AddLayer(layers[i-1], debug.GreenColor)
			debug.AddLayer(layers[i], debug.GreenColor)
			curP := calculateCenter(layers[i])
			//fmt.Printf("\n %v \n", layers[i])
			fmt.Printf("\n %v \n", i)

			//SaveSTL("/home/l1va/simp_mesh_before.stl", simpMesh)
			//debug.TriangleToDebugFile(curP, curP.Shift(AxisX.MulScalar(10)), curP.Shift(AxisY.MulScalar(15)), "debug_simplified.txt")
			simpMesh, _, err = helpers.CutMesh(simpMesh, Plane{P: curP, N: AxisZ})
			if err != nil {

				log.Fatal("failed to cut simpMesh by plane in 5a slicing: ", err)
			}
			mesh, _, err = helpers.CutMesh(mesh, Plane{P: curP, N: AxisZ})
			if err != nil {
				log.Fatal("failed to cut mesh by plane in 5a slicing: ", err)
			}
			debug.AddTriangleByPoints(newPlane.P, newPlane.P.Shift(AxisX.ProjectOnPlane(*newPlane).MulScalar(10)),
				newPlane.P.Shift(AxisY.ProjectOnPlane(*newPlane).MulScalar(15)), debug.GreenColor)

			//newPlane.N.Rotate(RotationAroundZ(90))

			//SaveSTL("/home/l1va/simp_mesh_before_cut.stl", simpMesh)

			simpMesh, _, err = helpers.CutMesh(simpMesh, *newPlane)
			if err != nil {
				log.Fatal("failed to cut simpMesh by newPlane in 5a slicing: ", err)
			}
			if i > 10 {
				debug.AddTriangleByPoints(newPlane.P, newPlane.P.Shift(AxisX.ProjectOnPlane(*newPlane).MulScalar(10)),
					newPlane.P.Shift(AxisY.ProjectOnPlane(*newPlane).MulScalar(15)), debug.GreenColor)

				//newPlane.N.Rotate(RotationAroundZ(90))

				//SaveSTL("/home/l1va/simp_mesh_before_cut.stl", simpMesh)

			}
			mesh, down, err = helpers.CutMesh(mesh, *newPlane)
			if err != nil {
				log.Fatal("failed to cut mesh by newPlane in 5a slicing: ", err)
			}
			toAdd = SliceByVector(down, AxisZ, settings)
			gcd.Add(gcode.LayersMoving{Layers: toAdd, Index: gcd.LayersCount})

			if rotated {
				simpMesh = simpMesh.Rotate(RotationAroundX(-angleX), OriginPoint)
				simpMesh = simpMesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)
				mesh = mesh.Rotate(RotationAroundX(-angleX), OriginPoint)
				mesh = mesh.Rotate(RotationAroundZ(-angleZ), OriginPoint)
				gcd.Add(gcode.RotateXZ{})
				rotated = false
			} else {
				angleZ = newPlane.N.ProjectOnPlane(PlaneXY).Angle(AxisX) + 90
				print(angleZ)
				simpMesh = simpMesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
				simpMesh = simpMesh.Rotate(RotationAroundX(angleX), OriginPoint)
				mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
				mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)
				gcd.Add(gcode.RotateXZ{AngleX: angleX, AngleZ: angleZ})
				rotated = true
			}
			layers = SliceByVector(simpMesh, AxisZ, settings)
			toAdd = SliceByVector(mesh, AxisZ, settings)
			//cmds = append(cmds, gcode.LayersMoving{layers, layersCount})
			//break
			i = 0
		}
	}

	return gcd
}

func CalculateFails(prevLayer, curLayer Layer) *Plane {
	rm := RotationAroundZ(float64(angleDist))

	for _, pth := range curLayer.Paths {
		curCp := FindCentroid(pth)
		if len(pth.Points) > 1 { //skip pth if it is just a hole
			p1 := pth.Points[0]
			p2 := pth.Points[1]
			if p1.VectorTo(p2).Cross(AxisZ).CodirectedWith(curCp.VectorTo(p1)) {
				continue
			}
		}

		//prevP := curCp.Shift(sh.Reverse())
		prevPath := findPrevPath(curCp, prevLayer)
		var failsCur []float64
		if prevPath == nil {
			println("prevPath not found")
			if len(prevLayer.Paths) > 0 {
				prevPath := prevLayer.Paths[0] //checking only first pth
				prevCp := FindCentroid(prevPath)
				v := prevCp.VectorTo(curCp).ProjectOnPlane(PlaneXY)
				pl := Plane{P: prevCp, N: v.Cross(AxisZ)}
				pi := pl.IntersectPathCodirectedWith(prevPath, v)
				if pi != nil {
					angleZ := prevCp.VectorTo(*pi).Angle(AxisX) + 90
					norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angleZ))
					return &Plane{P: *pi, N: norm}
				}
			}
		} else {
			curX := AxisX
			for j := 0; j < angleCount; j++ {
				pl := Plane{P: curCp, N: curX.Cross(AxisZ)}
				curPi := pl.IntersectPathCodirectedWith(pth, curX)
				if curPi != nil {
					prevPi := pl.IntersectPathCodirectedWith(*prevPath, curX)
					if prevPi != nil {
						if AxisZ.Angle(prevPi.VectorTo(*curPi)) > failAngle {
							failsCur = append(failsCur, float64(j*angleDist))
						}
					}
				}
				curX = curX.Rotate(rm)
			}
			if len(failsCur) > 0 {
				angleZ := float64(MiddleOnTheRing(failsCur, angleFull))
				v := AxisX.Rotate(RotationAroundZ(angleZ))
				pl := Plane{P: curCp, N: v.Cross(AxisZ)}

				pi := pl.IntersectPathCodirectedWith(*prevPath, v)
				if pi != nil {
					/*for _,li := range pth.Lines {
						pp := middlePoint(li.P1,li.P2)
						vv := li.P1.VectorTo(li.P2).Cross(AxisZ)
						debug.AddLine(Line{P1: pp, P2: pp.Shift(vv)})
					}*/
					norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angleZ + 90))
					//debug.AddLine(Line{P1: *pi, P2: pi.Shift(norm.MulScalar(10))})
					return &Plane{P: *pi, N: norm}
				}
			}
		}
	}
	return nil
}

func middlePoint(x1, x2 Point) Point {
	return Point{X: (x1.X + x2.X) / 2, Y: (x1.Y + x2.Y) / 2, Z: (x1.Z + x2.Z) / 2}
}

func findPrevPath(prevP Point, prevLayer Layer) *Path {
	for _, pth := range prevLayer.Paths {
		curP := FindCentroid(pth)
		if len(pth.Points) > 1 { //skip pth if it is just a hole
			p1 := pth.Points[0]
			p2 := pth.Points[1]
			if p1.VectorTo(p2).Cross(AxisZ).CodirectedWith(curP.VectorTo(p1)) {
				continue
			}
		}
		if prevP.Inside(pth) {
			return &pth
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
