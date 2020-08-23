package vip

import (
	"github.com/l1va/goosli/gcode"
	"github.com/l1va/goosli/helpers"
	. "github.com/l1va/goosli/primitives"
	"github.com/l1va/goosli/slicers"
	"log"
)

type BedPlane struct {
	tilted bool
	rotz   float64
}

func SliceByPlanes(mesh *Mesh, settings slicers.Settings, cutPlanes []AnalyzedPlane) gcode.Gcode {
	gcd := gcode.NewGcode(*settings.GcodeSettings)

	fillPlanes, fullFillPlanes := slicers.CalcFillPlanes(mesh, settings)
	var down *Mesh
	var err error

	bedPlanes := calcBedPlanes(cutPlanes)
	rotated := false
	for i := 0; i < len(bedPlanes); i++ {
		pl := bedPlanes[i]

		if i == len(cutPlanes) {
			down = mesh
		} else {
			mesh, down, err = helpers.CutMesh(mesh, cutPlanes[i].Plane)
			if err != nil {
				log.Fatal("failed to cut mesh, by plane: ", err, pl)
			}
		}

		if i != 0 {
			down = down.Rotate(RotationAroundZ(pl.rotz), settings.RotationCenter)
			gcd.Add(gcode.RotateZ{Angle: pl.rotz})
		}
		if pl.tilted {
			down = down.Rotate(RotationAroundX(-angleX), settings.RotationCenter)
			gcd.Add(gcode.InclineX{})
			rotated = true
		}
		if rotated && !pl.tilted {
			gcd.Add(gcode.InclineXBack{})
			rotated = false
		}
		add := slicers.SliceByVector(down, AxisZ, settings)
		gcd.AddLayers(PrepareLayers(add, settings, fillPlanes, fullFillPlanes))
		//TODO: fillPlanes fix
	}

	return gcd
}

func calcBedPlanes(cutPlanes []AnalyzedPlane) []BedPlane {
	res := []BedPlane{BedPlane{tilted: false, rotz: 0}}
	for _, p := range cutPlanes {
		res = append(res, BedPlane{tilted: p.tilted, rotz: p.rotz})
	}
	return res
}
