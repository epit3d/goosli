package vip

import (
	"github.com/l1va/goosli/gcode"
	. "github.com/l1va/goosli/primitives"
	. "github.com/l1va/goosli/slicers"
)

//1. Наклон стола будет осуществляться от команды М43, возврат в исходное положение М42
//2. Поворот стола не имеет определенного G-кода, а прописывается как четвертая координата.
// Т.е. команда выглядит как перемещение в точку G0 A236 B45, где А - координата поворота стола,
// 236 - фактическое значение угла поворота

var (
	angleX     = 60.0
	failAngle  = 40.0
	radiusDiff = 0.5
)

// Slice - Slicing on layers by simple algo
func Slice(mesh *Mesh, settings Settings) gcode.Gcode {

	planes := readPlanes(settings.PlanesFile)
	if len(planes)>0{
		return SliceByPlanes(mesh,settings, planes)
	}

	layers := SliceByVector(mesh, settings.LayerHeight, AxisZ)

	op := getAllOutsidePathes(layers[0])
	if len(op) == 2 {
		return SliceBridge(mesh, settings, layers)
	}
	println("outside pathes: ", len(op), len(op[0].Lines))
	if len(op) == 1 && isRotation(op[0]) {
		op = getAllOutsidePathes(layers[len(layers)-1])
		println("outside pathes: ", len(op), len(op[0].Lines))
		if len(op) == 1 && isRotation(op[0]) {
			return SliceRotation(mesh, settings, layers)
		}
	}
	return SliceDefault2(mesh, settings, layers)
}

func isRotation(pth Path) bool {
	cp := FindCentroid(pth)
	d := cp.DistanceTo(pth.Lines[0].P1)
	if len(pth.Lines) < 14 { //TODO: hardcode! square is a circle too but with small count of points
		return false
	}
	for _, line := range pth.Lines {
		if cp.DistanceTo(line.P2) < d-radiusDiff || cp.DistanceTo(line.P2) > d+radiusDiff {
			return false
		}
	}
	return true
}
