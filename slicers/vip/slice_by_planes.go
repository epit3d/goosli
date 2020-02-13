package vip

import (
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers"
	"log"
)

func SliceByPlanes(mesh *Mesh, settings slicers.Settings, planes []AnalyzedPlane) gcode.Gcode {
	var gcd gcode.Gcode
	fillPlanes := slicers.CalcFillPlanes(mesh, settings)
	var down *Mesh
	var err error
	firstPlane := AnalyzedPlane{tilted: false, rotz: 0, Plane: PlaneXY}
	planes = append([]AnalyzedPlane{firstPlane}, planes...)
	rotated := false
	for i := 0; i < len(planes); i++ {
		pl := planes[i]

		if i == len(planes)-1 {
			down = mesh
		} else {
			mesh, down, err = helpers.CutMesh(mesh, planes[i+1].Plane)
			if err != nil {
				log.Fatal("failed to cut mesh, by plane: ", err, pl)
			}
		}

		if i != 0 {
			down = down.Rotate(RotationAroundZ(pl.rotz+180), OriginPoint)
			gcd.Add(gcode.RotateZ{Angle: pl.rotz + 180})
		}
		if pl.tilted {
			down = down.Rotate(RotationAroundX(angleX), OriginPoint)
			gcd.Add(gcode.InclineX{})
			rotated = true
		}
		if rotated && !pl.tilted {
			gcd.Add(gcode.InclineXBack{})
			rotated = false
		}
		add := slicers.SliceByVector(down, settings.LayerHeight, AxisZ)
		gcd.Add(gcode.LayersMoving{Layers: PrepareLayers(add, settings, fillPlanes), Index: gcd.LayersCount})
		//TODO: fillPlanes fix
	}

	return gcd
}
