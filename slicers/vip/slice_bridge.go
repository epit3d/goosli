package vip

import (
	"github.com/l1va/goosli/debug"
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	. "github.com/l1va/goosli/slicers"
	"log"
)

func SliceBridge(mesh *Mesh, settings Settings, layers []Layer) gcode.Gcode {
	debug.RecreateFile()
	println("BRIDGE slicing:")
	var gcd gcode.Gcode

	fillPlanes := CalcFillPlanes(mesh, settings)

	i := 0
	var newPlane *Plane
	for i < len(layers) {
		if i == 0 {
			i++
			continue
		}
		newPlane = failExistBridge(layers[i-1], layers[i])
		if newPlane != nil {
			println("LEVEL: ", i)
			break
		}

		i++
	}
	// Add horizontal layers
	gcd.Add(gcode.LayersMoving{PrepareLayers(layers[:i], settings, fillPlanes), gcd.LayersCount, settings.GetExtrusionParams()})

	anyPoint := layers[i].Paths[0].Points[0]
	mesh, down, err := helpers.CutMesh(mesh, Plane{anyPoint, AxisZ})
	if err != nil {
		log.Fatal("failed to cut mesh, before rotation: ", err)
	}
	// Add horizontal layers for one leg before rotation
	mesh, down, err = helpers.CutMesh(mesh, *newPlane)
	if err != nil {
		log.Fatal("failed to cut mesh, for one leg: ", err)
	}
	oneLeg := SliceByVector(down, AxisZ, settings)
	gcd.Add(gcode.LayersMoving{PrepareLayers(oneLeg, settings, fillPlanes), gcd.LayersCount, settings.GetExtrusionParams()})

	// Rotate bed
	angleZ := newPlane.N.ProjectOnPlane(PlaneXY).Angle(AxisX) + 90
	mesh = mesh.Rotate(RotationAroundZ(angleZ), OriginPoint)
	mesh = mesh.Rotate(RotationAroundX(angleX), OriginPoint)
	gcd.Add(gcode.RotateZ{Angle: angleZ})
	gcd.Add(gcode.InclineX{})

	for i, plane := range fillPlanes {
		fillPlanes[i] = plane.Rotate(RotationAroundZ(angleZ)).Rotate(RotationAroundX(angleX))
	}

	rest := SliceByVector(mesh, AxisZ, settings)
	gcd.Add(gcode.LayersMoving{PrepareLayers(rest, settings, fillPlanes), gcd.LayersCount, settings.GetExtrusionParams()})

	// Rotate bed back
	gcd.Add(gcode.InclineXBack{})
	gcd.Add(gcode.RotateZ{Angle: 0})
	return gcd
}

func failExistBridge(prevLayer, curLayer Layer) *Plane {

	prevP1, prevP2 := getTwoOutsidePathes(prevLayer)
	p1, p2 := getTwoOutsidePathes(curLayer)
	prevCp1 := FindCentroid(prevP1)
	cp1 := FindCentroid(p1)
	cp2 := FindCentroid(p2)

	if prevCp1.DistanceTo(cp1) > prevCp1.DistanceTo(cp2) {
		p1, p2 = p2, p1
		cp1, cp2 = cp2, cp1
	}

	// check fail for first leg
	plane := Plane{cp1, cp1.VectorTo(cp2).Cross(AxisZ)}

	ip1 := plane.IntersectPathCodirectedWith(p1, cp1.VectorTo(cp2))
	prevIp1 := plane.IntersectPathCodirectedWith(prevP1, cp1.VectorTo(cp2))
	if ip1 == nil || prevIp1 == nil {
		log.Println("could not find path intersections with plane: ", ip1, prevIp1, plane, p1, prevP1)
	} else if AxisZ.Angle(prevIp1.VectorTo(*ip1)) > failAngle {
		println(AxisZ.Angle(prevIp1.VectorTo(*ip1)))
		angle := AxisX.Angle(cp1.VectorTo(cp2))
		norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angle + 90))
		pl := Plane{*prevIp1, norm}
		return &pl
	}

	// check fail for second leg
	ip2 := plane.IntersectPathCodirectedWith(p2, cp2.VectorTo(cp1))
	prevIp2 := plane.IntersectPathCodirectedWith(prevP2, cp2.VectorTo(cp1))
	if ip2 == nil || prevIp2 == nil {
		log.Println("could not find path intersections with plane: ", ip2, prevIp2, plane, p2, prevP2)
	} else if AxisZ.Angle(prevIp2.VectorTo(*ip2)) > failAngle {
		println(AxisZ.Angle(prevIp2.VectorTo(*ip2)))
		angle := AxisX.Angle(cp2.VectorTo(cp1))
		norm := AxisZ.Rotate(RotationAroundX(angleX)).Rotate(RotationAroundZ(angle + 90))
		pl := Plane{*prevIp2, norm}
		return &pl
	}

	return nil
}
func getTwoOutsidePathes(layer Layer) (Path, Path) {
	all := getAllOutsidePathes(layer)
	if len(all) == 2 {
		return all[0], all[1]
	}
	//debug.AddLayer(layer) //TODO:
	println("there are not two outside pathes: ", len(all))
	return all[0], all[1]
}
func getAllOutsidePathes(layer Layer) []Path { //for horizontal pathes
	var res []Path
	for i, pth := range layer.Paths {
		curP := FindCentroid(pth)

		if len(pth.Points) > 1 { //skip path if it is just an inner hole
			p1 := pth.Points[0]
			p2 := pth.Points[1]
			if p1.VectorTo(p2).Cross(AxisZ).CodirectedWith(curP.VectorTo(p1)) {
				continue
			}
		}
		res = append(res, layer.Paths[i])
	}
	return res
}
